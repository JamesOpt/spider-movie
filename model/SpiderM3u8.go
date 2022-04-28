package model

type SpiderM3u8 struct {
	Model
	SeriesId  uint `gorm:"series_id"`
	Status int
	Link string
	TargetLink string
	Filename string
}