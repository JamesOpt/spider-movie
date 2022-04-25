package collector

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"spider-movie/hleper"
)

var DOMAIN = "https://www.pianba.tv"

type PianBa struct {
	request *hleper.Request
	
}

func (p *PianBa) Run(path string)  {
	c := colly.NewCollector()

	c.OnResponse(func(response *colly.Response) {
		reader := bytes.NewReader(response.Body)
		dom, err := goquery.NewDocumentFromReader(reader)
		if err != nil{
			panic(err)
		}

		// 获取类型
		dom.Find("ul.stui-header__menu .active").Text()
	})
	
	//c.OnHTML("div.stui-content", func(e *colly.HTMLElement) {
	//	fmt.Printf(e.ChildAttr(".stui-content__thumb img", "data-original"))
	//	fmt.Printf(e.ChildText(".stui-content__detail h1:first-child"))
	//	fmt.Printf(e.ChildText(".stui-content__detail h1:first-child"))
	//})

	c.Visit("https://www.pianba.tv/html/194890.html")
}