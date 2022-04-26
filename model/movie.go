package model

type Movie struct {
	Model
	Title       string
	Episodes    int64
	SpiderLink  string
	Duration    string
	Douban      string
	Cover       string
	Lang        string
	Type        int
	Description string
	Genre       string
	EnAlias     string   `gorm:"en_alias"`
	Serises     []Series `gorm:foreignKey:movie_id`
}