package repo

import (
	model "permen_api/domain/sample/model"
)

const (
	GetAllUserIntegrationsQuery = `SELECT username, credentials, created_by, channel_name, is_active, created_at, updated_at FROM user_integration`
)

func (r *userIntegrationRepo) GetAllUserIntegrations() ([]*model.UserIntegration, error) {
	var userIntegrations []*model.UserIntegration
	err := r.db.Raw(GetAllUserIntegrationsQuery).Scan(&userIntegrations).Error
	if err != nil {
		return nil, err
	}
	return userIntegrations, nil
}
