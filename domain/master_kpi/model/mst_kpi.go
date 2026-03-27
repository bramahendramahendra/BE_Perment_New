package model

type (
	MstPerspektif struct {
		IdPerspektif int    `gorm:"column:id_perspektif"`
		Perspektif   string `gorm:"column:perspektif"`
	}
)
