package collector

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/url"
	"regexp"
	"spider-movie/db"
	"spider-movie/helper"
	"spider-movie/model"
	"strings"
)

var DOMAIN = "https://www.pianba.tv"

type PianBa struct {
	Request *helper.Request
	movie model.Movie
}

func (p *PianBa) Run(u string)  {
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		p.movie.SpiderLink = request.URL.String()
	})

	c.OnResponse(func(response *colly.Response) {
		reader := bytes.NewReader(response.Body)
		dom, err := goquery.NewDocumentFromReader(reader)
		if err != nil{
			panic(err)
		}

		// 获取类型
		if "电影" == dom.Find("ul.stui-header__menu .active").Text() {
			p.movie.Type = 1
		}

		p.movie.Cover, _ = dom.Find(".stui-content__thumb img").Attr("data-original")
		p.movie.Title = dom.Find(".stui-content__detail .title").Text()
		db.Engine.Driver().Where("title = ?", p.movie.Title).FirstOrCreate(&p.movie)
	})
	
	//c.OnHTML("div.stui-content", func(e *colly.HTMLElement) {
	//	p.movie.Cover = e.ChildAttr(".stui-content__thumb img", "data-original")
	//	p.movie.Title = e.ChildText(".stui-content__detail h1:first-child")
	//})

	c.OnHTML(".stui-content__playlist", func(element *colly.HTMLElement) {
		// 只搜索第一个dom
		if element.Index == 0{
			for serial, link := range element.ChildAttrs("a", "href") {

				p.Request.SetUrl(DOMAIN, link, nil)
				response := p.Request.Get()

				m3u8FileLink := regexp.MustCompile(`http.*?\.m3u8`).FindString(response)
				m3u8FileLink = strings.Replace(m3u8FileLink, "\\", "", -1)
				uu, _ := url.Parse(m3u8FileLink)

				p.Request.SetUrl(uu.Scheme + "://" + uu.Host, uu.Path, uu.Query())

				realLink := p.Request.Get()
				realLink = regexp.MustCompile(`.+\.m3u8`).FindString(realLink)
				hosts := regexp.MustCompile(`http://|https://[^/]+`).FindString(m3u8FileLink)
				realLink = hosts + realLink

				serialModel := &model.Series{}
				db.Engine.Driver().FirstOrCreate(serialModel, model.Series{
					MovieId: p.movie.ID,
					Serial: serial + 1,
				})

				fmt.Printf("第%v集, %v\n", serial + 1, realLink)
				DownloadRaw(realLink, p.movie.Title, serialModel)
			}
		}
	})

	c.Visit(u)
}