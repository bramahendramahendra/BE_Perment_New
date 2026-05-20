package edm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDBQuerier implementasi DBQuerier untuk testing.
type mockDBQuerier struct {
	scanResult any
	scanErr    error
}

func (m *mockDBQuerier) RawScan(_ string, dest any, _ ...any) error {
	if m.scanErr != nil {
		return m.scanErr
	}
	if m.scanResult != nil {
		b, _ := json.Marshal(m.scanResult)
		_ = json.Unmarshal(b, dest)
	}
	return nil
}

func newTestClient(db DBQuerier, serverURL string) *edmClient {
	return &edmClient{
		db:         db,
		httpClient: &http.Client{},
		debug:      false,
	}
}

func TestGetKpi_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, datahubChannel, r.Header.Get("X-DATAHUB-CHANNEL"))
		assert.Equal(t, datahubPersonalNumber, r.Header.Get("X-DATAHUB-PERSONAL-NUMBER"))
		assert.NotEmpty(t, r.Header.Get("Authorization"))

		resp := customerhubResponse{
			StatusCode:      200,
			ErrorCode:       "000",
			ResponseCode:    "00",
			ResponseMessage: "Success",
			Data: []KpiItem{
				{ID: "014002002001002", AliasName: "Pendapatan Bunga", Amount: 17550444967},
				{ID: "014002003003000", AliasName: "%Rate FTP KPR", Amount: 0.541},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	db := &mockDBQuerier{
		scanResult: paramRow{
			Userid:   "cekadmin",
			Userpass: "cekadmin",
			Userurl:  server.URL,
		},
	}

	client := newTestClient(db, server.URL)
	result, err := client.GetKpi("2025-02")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "014002002001002", result[0].ID)
	assert.Equal(t, "Pendapatan Bunga", result[0].AliasName)
	assert.InDelta(t, 17550444967.0, result[0].Amount, 0.01)
	assert.Equal(t, "014002003003000", result[1].ID)
	assert.InDelta(t, 0.541, result[1].Amount, 0.001)
}

func TestGetKpi_DBError(t *testing.T) {
	db := &mockDBQuerier{
		scanErr: assert.AnError,
	}

	client := newTestClient(db, "")
	result, err := client.GetKpi("2025-02")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal ambil param")
	assert.Nil(t, result)
}

func TestGetKpi_APIResponseCodeFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := customerhubResponse{
			StatusCode:      400,
			ErrorCode:       "001",
			ResponseCode:    "01",
			ResponseMessage: "Bad Request",
			Data:            nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	db := &mockDBQuerier{
		scanResult: paramRow{
			Userid:   "cekadmin",
			Userpass: "cekadmin",
			Userurl:  server.URL,
		},
	}

	client := newTestClient(db, server.URL)
	result, err := client.GetKpi("2025-02")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Customerhub GetKpi gagal")
	assert.Nil(t, result)
}

func TestGetKpi_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	db := &mockDBQuerier{
		scanResult: paramRow{
			Userid:   "cekadmin",
			Userpass: "cekadmin",
			Userurl:  server.URL,
		},
	}

	client := newTestClient(db, server.URL)
	result, err := client.GetKpi("2025-02")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal decode response")
	assert.Nil(t, result)
}

func TestBasicAuth(t *testing.T) {
	encoded := basicAuth("cekadmin", "cekadmin")
	assert.Equal(t, "Y2VrYWRtaW46Y2VrYWRtaW4=", encoded)
}

func TestGetKpi_EmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := customerhubResponse{
			StatusCode:      200,
			ResponseCode:    "00",
			ResponseMessage: "Success",
			Data:            []KpiItem{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	db := &mockDBQuerier{
		scanResult: paramRow{
			Userid:   "cekadmin",
			Userpass: "cekadmin",
			Userurl:  server.URL,
		},
	}

	client := newTestClient(db, server.URL)
	result, err := client.GetKpi("2025-02")

	require.NoError(t, err)
	assert.Empty(t, result)
}
