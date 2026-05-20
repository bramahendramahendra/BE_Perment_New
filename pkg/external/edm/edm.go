package edm

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"permen_api/config"
	"time"

	"gorm.io/gorm"
)

const (
	vendorGetKpi   = "GetKpi"
	defaultTimeout = 30 * time.Second

	datahubChannel        = "DATAHUB"
	datahubPersonalNumber = "00000000"
)

// DBQuerier abstraksi operasi DB yang dibutuhkan EdmClient.
type DBQuerier interface {
	RawScan(query string, dest any, args ...any) error
}

type (
	KpiItem struct {
		ID        string  `json:"id"`
		AliasName string  `json:"aliasName"`
		Amount    float64 `json:"amount"`
	}

	EdmClient interface {
		GetKpi(periode string) ([]KpiItem, error)
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

	customerhubResponse struct {
		StatusCode      int       `json:"statusCode"`
		ErrorCode       string    `json:"errorCode"`
		ResponseCode    string    `json:"responseCode"`
		ResponseMessage string    `json:"responseMessage"`
		Data            []KpiItem `json:"data"`
	}
)

// gormQuerier adalah adapter DBQuerier untuk *gorm.DB.
type gormQuerier struct {
	db *gorm.DB
}

func (g *gormQuerier) RawScan(query string, dest any, args ...any) error {
	return g.db.Raw(query, args...).Scan(dest).Error
}

func New(db *gorm.DB, debug bool) EdmClient {
	return &edmClient{
		db: &gormQuerier{db: db},
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.ENV.ReleaseMode != "production"},
			},
		},
		debug: debug,
	}
}

// GetKpi mengambil data KPI dari Customerhub API berdasarkan periode.
func (c *edmClient) GetKpi(periode string) ([]KpiItem, error) {
	param, err := c.getParam(vendorGetKpi)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil param %s: %w", vendorGetKpi, err)
	}

	body := map[string]any{
		"periode": periode,
	}

	result, err := c.post(param.Userurl, param.Userid, param.Userpass, body)
	if err != nil {
		return nil, fmt.Errorf("gagal request GetKpi ke Customerhub: %w", err)
	}

	if result.ResponseCode != "00" {
		return nil, fmt.Errorf("Customerhub GetKpi gagal: %s", result.ResponseMessage)
	}

	return result.Data, nil
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

// post melakukan HTTP POST ke url dengan JSON body dan Basic Auth.
func (c *edmClient) post(url, username, password string, body map[string]any) (*customerhubResponse, error) {
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DATAHUB-CHANNEL", datahubChannel)
	req.Header.Set("X-DATAHUB-PERSONAL-NUMBER", datahubPersonalNumber)
	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))

	if c.debug {
		fmt.Printf("[DEBUG] Customerhub POST %s\n", url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result customerhubResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response Customerhub: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Customerhub response status: %d\n", resp.StatusCode)
	}

	return &result, nil
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
