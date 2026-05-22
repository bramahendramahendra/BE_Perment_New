package model

type (
	MstLinkFormat struct {
		IdLinkFormat int    `gorm:"column:id_link_format"`
		UrlPrefix    string `gorm:"column:url_prefix"`
		Keterangan   string `gorm:"column:keterangan"`
	}
)
