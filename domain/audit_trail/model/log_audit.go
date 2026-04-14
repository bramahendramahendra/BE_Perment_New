package model

import "time"

type LogAudit struct {
	Ipaddress   string    `gorm:"column:ipaddress"`
	Userid      string    `gorm:"column:userid"`
	Function    string    `gorm:"column:function"`
	Body        string    `gorm:"column:body"`
	Response    string    `gorm:"column:response"`
	Errordesc   string    `gorm:"column:errordesc"`
	TimeRequest time.Time `gorm:"column:time_request"`
}

func (LogAudit) TableName() string {
	return "log_request"
}
