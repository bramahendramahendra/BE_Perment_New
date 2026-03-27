package dto

type (
	MasterPerspektifResponse struct {
		IdPerspektif int    `gorm:"column:id_perspektif"`
		Perspektif   string `gorm:"column:perspektif"`
	}
)
