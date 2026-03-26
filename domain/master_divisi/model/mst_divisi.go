package model

type MstDivisi struct {
	Kostl   string `gorm:"column:KOSTL"`
	KostlTx string `gorm:"column:KOSTL_TX"`
}
