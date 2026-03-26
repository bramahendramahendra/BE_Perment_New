package dto

type (
	MasterDivisiResponse struct {
		Kostl   string `gorm:"column:KOSTL"    json:"kostl"`
		KostlTx string `gorm:"column:KOSTL_TX" json:"kostlTx"`
	}
)
