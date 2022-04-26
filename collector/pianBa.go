package collector

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/url"
	"os"
	"regexp"
	"spider-movie/hleper"
	"strings"
)

var DOMAIN = "https://www.pianba.tv"

type PianBa struct {
	Request *hleper.Request
	
}

func (p *PianBa) Run(u string)  {
	c := colly.NewCollector()

	c.OnResponse(func(response *colly.Response) {
		reader := bytes.NewReader(response.Body)
		dom, err := goquery.NewDocumentFromReader(reader)
		if err != nil{
			panic(err)
		}

		// 获取类型
		fmt.Println(dom.Find("ul.stui-header__menu .active").Text())
	})
	
	c.OnHTML("div.stui-content", func(e *colly.HTMLElement) {
		//fmt.Printf(e.ChildAttr(".stui-content__thumb img", "data-original"))
		//fmt.Printf(e.ChildText(".stui-content__detail h1:first-child"))
		//fmt.Printf(e.ChildText(".stui-content__detail h1:first-child"))
	})

	c.OnHTML(".stui-content__playlist", func(element *colly.HTMLElement) {

		if element.Index == 0{
			for _, link := range element.ChildAttrs("a", "href") {
				//fmt.Println(DOMAIN + link)
				p.Request.Url = hleper.Url{
					Host:  DOMAIN,
					Path:  link,
				}
				response := p.Request.Get()

				m3u8FileLink := regexp.MustCompile(`http.*?\.m3u8`).FindString(response)
				m3u8FileLink = strings.Replace(m3u8FileLink, "\\", "", -1)
				uu, _ := url.Parse(m3u8FileLink)

				p.Request.SetUrl(uu.Scheme + "://" + uu.Host, uu.Path, nil)

				realLink := p.Request.Get()
				realLink = regexp.MustCompile(`.+\.m3u8`).FindString(realLink)
				hosts := regexp.MustCompile(`http://|https://[^/]+`).FindString(m3u8FileLink)
				realLink = hosts + realLink
				uu, _ = url.Parse(realLink)

				p.Request.SetUrl(uu.Scheme + "://" + uu.Host, uu.Path, nil)

				//pp,_ := os.Getwd()
				m := M3u8{}
				m.DownloadRaw(realLink)
				//p.Request.Download("test.m3u8", filepath.Join(pp, ".."))
os.Exit(1)
			}
		}


	})

	c.Visit(u)
}