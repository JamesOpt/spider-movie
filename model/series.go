package model

type Series struct {
	Model
	MovieId    uint `gorm:"movie_id"`
	Status     int64
	Serial     int
	LocalLink  string
	SpiderLink string
	OtherLink  string
}
