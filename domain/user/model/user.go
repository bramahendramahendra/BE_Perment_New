package model

// User merepresentasikan struktur data dari tabel `user`
type User struct {
	PERNR        string `gorm:"column:PERNR"`
	SNAME        string `gorm:"column:SNAME"`
	JGPG         string `gorm:"column:JGPG"`
	ESELON       string `gorm:"column:ESELON"`
	WERKS        string `gorm:"column:WERKS"`
	WERKS_TX     string `gorm:"column:WERKS_TX"`
	BTRTL        string `gorm:"column:BTRTL"`
	BTRTL_TX     string `gorm:"column:BTRTL_TX"`
	KOSTL        string `gorm:"column:KOSTL"`
	KOSTL_TX     string `gorm:"column:KOSTL_TX"`
	ORGEH        string `gorm:"column:ORGEH"`
	ORGEH_TX     string `gorm:"column:ORGEH_TX"`
	STELL        string `gorm:"column:STELL"`
	STELL_TX     string `gorm:"column:STELL_TX"`
	PLANS        string `gorm:"column:PLANS"`
	PLANS_TX     string `gorm:"column:PLANS_TX"`
	HILFM        string `gorm:"column:HILFM"`
	HTEXT        string `gorm:"column:HTEXT"`
	BRANCH       string `gorm:"column:BRANCH"`
	MAINBR       string `gorm:"column:MAINBR"`
	IS_PEMIMPIN  string `gorm:"column:IS_PEMIMPIN"`
	ADMIN_LEVEL  string `gorm:"column:ADMIN_LEVEL"`
	ORGEH_PGS    string `gorm:"column:ORGEH_PGS"`
	ORGEH_PGS_TX string `gorm:"column:ORGEH_PGS_TX"`
	PLANS_PGS    string `gorm:"column:PLANS_PGS"`
	PLANS_PGS_TX string `gorm:"column:PLANS_PGS_TX"`
	BRANCH_PGS   string `gorm:"column:BRANCH_PGS"`
	HILFM_PGS    string `gorm:"column:HILFM_PGS"`
	HTEXT_PGS    string `gorm:"column:HTEXT_PGS"`
	TIPE_UKER    string `gorm:"column:TIPE_UKER"`
	REKENING     string `gorm:"column:REKENING"`
	NPWP         string `gorm:"column:NPWP"`
	REGION       string `gorm:"column:REGION"`
	RGDESC       string `gorm:"column:RGDESC"`
	BRDESC       string `gorm:"column:BRDESC"`
	MBDESC       string `gorm:"column:MBDESC"`
}
