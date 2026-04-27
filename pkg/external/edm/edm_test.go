package edm

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// Mock DBQuerier
// =============================================================================

type mockDB struct {
	// rawResults memetakan query prefix ke nilai yang akan di-scan ke dest.
	rawResults map[string]interface{}
	rawErr     error
	execErr    error

	capturedExecQuery string
	capturedExecArgs  []interface{}
}

func (m *mockDB) RawScan(query string, dest interface{}, args ...interface{}) error {
	if m.rawErr != nil {
		return m.rawErr
	}
	val, ok := m.rawResults[query]
	if !ok {
		return nil
	}
	switch d := dest.(type) {
	case *int64:
		if v, ok := val.(int64); ok {
			*d = v
		}
	case *string:
		if v, ok := val.(string); ok {
			*d = v
		}
	case *paramRow:
		if v, ok := val.(paramRow); ok {
			*d = v
		}
	}
	return nil
}

func (m *mockDB) Exec(query string, args ...interface{}) error {
	m.capturedExecQuery = query
	m.capturedExecArgs = args
	return m.execErr
}

// =============================================================================
// Helper: buat edmClient dengan mock DB dan mock HTTP server
// =============================================================================

func newTestClient(db DBQuerier, httpClient *http.Client) *edmClient {
	return &edmClient{
		db:         db,
		httpClient: httpClient,
		debug:      false,
	}
}

func newMockServer(responseBody interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(responseBody)
	}))
}

// =============================================================================
// TestGetToken
// =============================================================================

func TestGetToken_Success(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"access_token": "token-abc123",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid:   "client-id",
				Userpass: "client-secret",
				Userurl:  srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	token, err := c.GetToken()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token != "token-abc123" {
		t.Errorf("expected token 'token-abc123', got: %s", token)
	}
	if db.capturedExecQuery == "" {
		t.Error("expected Exec dipanggil untuk menyimpan token, tapi tidak dipanggil")
	}
}

func TestGetToken_ParamDBError(t *testing.T) {
	db := &mockDB{rawErr: errors.New("db connection refused")}
	c := newTestClient(db, http.DefaultClient)

	_, err := c.GetToken()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, db.rawErr) {
		t.Errorf("expected error wrapping db error, got: %v", err)
	}
}

func TestGetToken_EDMReturnError(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"error": "invalid_client",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	_, err := c.GetToken()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetToken_MissingAccessToken(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"message": "ok",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	_, err := c.GetToken()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetToken_SaveToDB_Error(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"access_token": "token-xyz",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
		execErr: errors.New("db write error"),
	}

	c := newTestClient(db, srv.Client())
	_, err := c.GetToken()

	if err == nil {
		t.Fatal("expected error saat simpan token, got nil")
	}
}

// =============================================================================
// TestGetDataKPI
// =============================================================================

func TestGetDataKPI_Success(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"success": true,
		"data": []interface{}{
			map[string]interface{}{"id_kpi": "KPI-001", "nilai": 95.5},
		},
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			// token masih valid (count = 0 expired)
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(0),
			// token dari cache
			"SELECT token FROM param_token_edm LIMIT 1": "cached-token",
			// param GetDataKPI
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	data, err := c.GetDataKPI("2024", "1", "KPI-001")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if data == nil {
		t.Error("expected data, got nil")
	}
}

func TestGetDataKPI_SuccessFalse(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"success": false,
		"message": "data tidak ditemukan",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(0),
			"SELECT token FROM param_token_edm LIMIT 1": "cached-token",
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	_, err := c.GetDataKPI("2024", "1", "KPI-001")

	if err == nil {
		t.Fatal("expected error karena success=false, got nil")
	}
}

func TestGetDataKPI_EmptyData(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(0),
			"SELECT token FROM param_token_edm LIMIT 1": "cached-token",
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "x", Userpass: "y", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	_, err := c.GetDataKPI("2024", "1", "KPI-001")

	if err == nil {
		t.Fatal("expected error karena data kosong, got nil")
	}
}

func TestGetDataKPI_TokenExpired_RefreshSuccess(t *testing.T) {
	// Server untuk GetToken
	tokenSrv := newMockServer(map[string]interface{}{
		"access_token": "new-token",
	}, http.StatusOK)
	defer tokenSrv.Close()

	// Server untuk GetDataKPI
	dataSrv := newMockServer(map[string]interface{}{
		"success": true,
		"data":    []interface{}{map[string]interface{}{"id": "KPI-001"}},
	}, http.StatusOK)
	defer dataSrv.Close()

	callCount := 0
	db := &mockDB{
		rawResults: map[string]interface{}{
			// token expired
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(1),
		},
	}

	// Override RawScan agar param bisa dibedakan antara GetToken dan GetDataKPI
	db.rawResults["SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?"] = paramRow{}

	// Kita pakai custom mock yang bisa membedakan urutan panggilan
	customDB := &sequentialMockDB{
		paramResponses: []paramRow{
			{Userid: "id", Userpass: "secret", Userurl: tokenSrv.URL}, // GetToken
			{Userid: "id", Userpass: "secret", Userurl: dataSrv.URL},  // GetDataKPI
		},
		expiredCount: int64(1),
	}
	_ = callCount

	c := newTestClient(customDB, tokenSrv.Client())
	c.httpClient = &http.Client{} // pakai default agar bisa hit keduanya

	data, err := c.GetDataKPI("2024", "1", "KPI-001")

	if err != nil {
		t.Fatalf("expected no error setelah token refresh, got: %v", err)
	}
	if data == nil {
		t.Error("expected data, got nil")
	}
}

// =============================================================================
// TestGetOrRefreshToken
// =============================================================================

func TestGetOrRefreshToken_UseCachedToken(t *testing.T) {
	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(0),
			"SELECT token FROM param_token_edm LIMIT 1": "cached-token-xyz",
		},
	}

	c := newTestClient(db, http.DefaultClient)
	token, err := c.getOrRefreshToken()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token != "cached-token-xyz" {
		t.Errorf("expected cached token, got: %s", token)
	}
}

func TestGetOrRefreshToken_EmptyCache_Refresh(t *testing.T) {
	srv := newMockServer(map[string]interface{}{
		"access_token": "refreshed-token",
	}, http.StatusOK)
	defer srv.Close()

	db := &mockDB{
		rawResults: map[string]interface{}{
			"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?": int64(0),
			"SELECT token FROM param_token_edm LIMIT 1": "",
			"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?": paramRow{
				Userid: "id", Userpass: "secret", Userurl: srv.URL,
			},
		},
	}

	c := newTestClient(db, srv.Client())
	token, err := c.getOrRefreshToken()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token != "refreshed-token" {
		t.Errorf("expected 'refreshed-token', got: %s", token)
	}
}

// =============================================================================
// sequentialMockDB — helper untuk test yang butuh urutan panggilan berbeda
// =============================================================================

type sequentialMockDB struct {
	paramResponses []paramRow
	paramCallIdx   int
	expiredCount   int64
	execErr        error
}

func (s *sequentialMockDB) RawScan(query string, dest interface{}, args ...interface{}) error {
	switch d := dest.(type) {
	case *int64:
		*d = s.expiredCount
	case *string:
		// token cache kosong — paksa refresh
		*d = ""
	case *paramRow:
		if s.paramCallIdx < len(s.paramResponses) {
			*d = s.paramResponses[s.paramCallIdx]
			s.paramCallIdx++
		}
	}
	return nil
}

func (s *sequentialMockDB) Exec(query string, args ...interface{}) error {
	return s.execErr
}
