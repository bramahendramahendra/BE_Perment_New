package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dto "permen_api/domain/edm/dto"
)

type edmTokenRow struct {
	Token      string
	InsertDate time.Time
}

type edmParamRow struct {
	Userid  string
	Userpass string
	Userurl string
}

func (s *edmService) getEDMParam(vendor string) (edmParamRow, error) {
	var row edmParamRow
	err := s.db.Raw("SELECT userid, userpass, userurl FROM mst_param WHERE vendor = ?", vendor).
		Scan(&row).Error
	return row, err
}

func (s *edmService) getOrRefreshToken() (string, error) {
	var count int64
	s.db.Raw("SELECT COUNT(*) FROM param_token_edm WHERE TIMESTAMPDIFF(HOUR, insert_date, NOW()) >= 11").Scan(&count)

	if count > 0 {
		return s.fetchAndStoreToken()
	}

	var token string
	s.db.Raw("SELECT token FROM param_token_edm LIMIT 1").Scan(&token)
	if token == "" {
		return s.fetchAndStoreToken()
	}
	return token, nil
}

func (s *edmService) fetchAndStoreToken() (string, error) {
	param, err := s.getEDMParam("GetToken")
	if err != nil {
		return "", fmt.Errorf("gagal ambil param GetToken: %w", err)
	}

	body := map[string]interface{}{
		"client_id":     param.Userid,
		"client_secret": param.Userpass,
		"grant_type":    "client_credentials",
	}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", param.Userurl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result["error"] != nil {
		return "", fmt.Errorf("EDM token error: %v", result["error"])
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("access_token tidak ditemukan dalam response EDM")
	}

	s.db.Exec("UPDATE param_token_edm SET token = ?, insert_date = NOW()", token)

	return token, nil
}

func (s *edmService) GetRealisasi(req *dto.GetRealisasiRequest) (interface{}, error) {
	token, err := s.getOrRefreshToken()
	if err != nil {
		return nil, err
	}

	param, err := s.getEDMParam("GetDataKPI")
	if err != nil {
		return nil, fmt.Errorf("gagal ambil param GetDataKPI: %w", err)
	}

	body := map[string]interface{}{
		"TAHUN":   req.Tahun,
		"KUARTAL": req.Triwulan,
		"ID_KPI":  req.IdKpi,
	}
	bodyBytes, _ := json.Marshal(body)

	httpReq, err := http.NewRequest("POST", param.Userurl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
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
