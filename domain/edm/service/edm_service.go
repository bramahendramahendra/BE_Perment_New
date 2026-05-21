package service

import (
	"context"
	"encoding/json"
	"fmt"
	dto "permen_api/domain/edm/dto"
	edm "permen_api/pkg/external/edm"
	"time"
)

const kpiCacheTTL = 12 * time.Hour

func (s *edmService) GetKpi(req *dto.GetKpiRequest) (interface{}, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("edm:kpi:%s", req.Periode)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Bytes()
		if err == nil {
			var result []edm.KpiItem
			if err := json.Unmarshal(cached, &result); err == nil {
				return result, nil
			}
		}
	}

	data, err := s.edm.GetKpi(req.Periode)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if b, err := json.Marshal(data); err == nil {
			s.redis.Set(ctx, cacheKey, b, kpiCacheTTL)
		}
	}

	return data, nil
}
