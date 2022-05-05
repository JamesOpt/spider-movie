package collector

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"regexp"
	"spider-movie/db"
	"spider-movie/helper"
	"spider-movie/model"
)

var HKTV = "https://www.hktv03.com"

var HKTV_SPIDER_PATH = "https://jx2api.5408h.cn/m3u8/"

type Hktv struct {
	Request *helper.Request
	model.Movie
}

func (hk *Hktv) Run(uri string)  {
	c := colly.NewCollector()

	c.OnResponse(func(response *colly.Response) {
		dom, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))

		if err != nil {
			panic(err)
		}

		hk.Movie.Cover,_ = dom.Find("img.lazyload").Attr("data-original")

		selector := dom.Find(".myui-content__detail")
		hk.Movie.Title = selector.Find(".title").Text()

		db.Engine.Driver().Where("title = ?", hk.Movie.Title).FirstOrCreate(&hk.Movie)
	})

	c.OnHTML(".myui-content__list", func(element *colly.HTMLElement) {
		if element.Index == 1 {
			for serial, link := range element.ChildAttrs(".btn-default", "href") {
				response := hk.Request.Get(HKTV + link)


				m3u8FileLink := regexp.MustCompile(`"url":"(.*?\.m3u8)"`).FindStringSubmatch(response)

				serialModel := &model.Series{}
				db.Engine.Driver().FirstOrCreate(serialModel, model.Series{
					MovieId: hk.Movie.ID,
					Serial: serial + 1,
				})

				DownloadRaw(HKTV_SPIDER_PATH + m3u8FileLink[1], hk.Movie.Title, serialModel)


			}
		}
	})

	c.Visit(uri)
}