package edm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const (
	vendorGetToken   = "GetToken"
	vendorGetDataKPI = "GetDataKPI"
	tokenTTLHours    = 11
	defaultTimeout   = 30 * time.Second
)

// DBQuerier abstraksi operasi DB yang dibutuhkan EdmClient.
// Memudahkan unit testing tanpa koneksi database nyata.
type DBQuerier interface {
	RawScan(query string, dest interface{}, args ...interface{}) error
	Exec(query string, args ...interface{}) error
}

type (
	EdmClient interface {
		GetToken() (string, error)
		GetDataKPI(tahun, triwulan, idKPI string) (interface{}, error)
	}

	edmClient struct {
		db         DBQuerier
		httpClient *http.Client
		debug      bool
	}

	paramRow struct {
		Userid   string
		Userpass string
		Userurl  string
	}
)

// gormQuerier adalah adapter DBQuerier untuk *gorm.DB.
type gormQuerier struct {
	db *gorm.DB
}

func (g *gormQuerier) RawScan(query string, dest interface{}, args ...interface{}) error {
	return g.db.Raw(query, args...).Scan(dest).Error
}

func (g *gormQuerier) Exec(query string, args ...interface{}) error {
	return g.db.Exec(query, args...).Error
}

func New(db *gorm.DB, debug bool) EdmClient {
	return &edmClient{
		db: &gormQuerier{db: db},
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		debug: debug,
	}
}

// GetToken mengambil token baru dari EDM dan menyimpannya ke param_token_edm.
func (c *edmClient) GetToken() (string, error) {
	param, err := c.getParam(vendorGetToken)
	if err != nil {
		return "", fmt.Errorf("gagal ambil param %s: %w", vendorGetToken, err)
	}

	body := map[string]interface{}{
		"client_id":     param.Userid,
		"client_secret": param.Userpass,
		"grant_type":    "client_credentials",
	}

	result, err := c.post(param.Userurl, "", body)
	if err != nil {
		return "", fmt.Errorf("gagal request GetToken ke EDM: %w", err)
	}

	if result["error"] != nil {
		return "", fmt.Errorf("EDM token error: %v", result["error"])
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("access_token tidak ditemukan dalam response EDM")
	}

	if err := c.db.Exec("UPDATE param_token_edm SET token = ?, insert_date = NOW()", token); err != nil {
		return "", fmt.Errorf("gagal simpan token EDM: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] EDM GetToken: token berhasil diperbarui\n")
	}

	return token, nil
}

// GetDataKPI mengambil data KPI dari EDM berdasarkan tahun, triwulan, dan ID KPI.
func (c *edmClient) GetDataKPI(tahun, triwulan, idKPI string) (interface{}, error) {
	token, err := c.getOrRefreshToken()
	if err != nil {
		return nil, err
	}

	param, err := c.getParam(vendorGetDataKPI)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil param %s: %w", vendorGetDataKPI, err)
	}

	body := map[string]interface{}{
		"TAHUN":   tahun,
		"KUARTAL": triwulan,
		"ID_KPI":  idKPI,
	}

	result, err := c.post(param.Userurl, token, body)
	if err != nil {
		return nil, fmt.Errorf("gagal request GetDataKPI ke EDM: %w", err)
	}

	success, _ := result["success"].(bool)
	if !success {
		return nil, fmt.Errorf("EDM GetDataKPI mengembalikan success=false")
	}

	dataArr, ok := result["data"].([]interface{})
	if !ok || len(dataArr) == 0 {
		return nil, fmt.Errorf("data EDM kosong atau tidak valid")
	}

	return dataArr[0], nil
}

// getOrRefreshToken mengecek usia token di DB; refresh jika sudah >= tokenTTLHours jam.
func (c *edmClient) getOrRefreshToken() (string, error) {
	var count int64
	c.db.RawScan(
		"SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= ?",
		&count,
		tokenTTLHours,
	)

	if count > 0 {
		return c.GetToken()
	}

	var token string
	c.db.RawScan("SELECT token FROM param_token_edm LIMIT 1", &token)
	if token == "" {
		return c.GetToken()
	}

	return token, nil
}

// getParam mengambil kredensial dan URL endpoint dari tabel mst_param berdasarkan vendor.
func (c *edmClient) getParam(vendor string) (paramRow, error) {
	var row paramRow
	err := c.db.RawScan(
		"SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?",
		&row,
		vendor,
	)
	return row, err
}

// post melakukan HTTP POST ke url dengan JSON body dan optional Bearer token.
func (c *edmClient) post(url, token string, body map[string]interface{}) (map[string]interface{}, error) {
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if c.debug {
		fmt.Printf("[DEBUG] EDM POST %s\n", url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response EDM: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] EDM response status: %d\n", resp.StatusCode)
	}

	return result, nil
}
