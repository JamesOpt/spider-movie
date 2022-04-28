package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"spider-movie/collector"
	"spider-movie/helper"
	"strings"
)

func main(){
	c := collector.PianBa{
		Request: helper.NewRequest(),
	}
	c.Run("https://www.pianba.tv/html/209094.html")
}

//func requestTest()  {
//	request := helper.NewRequest()
//	request.Header.Add("Authorization", "bearereyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9kZXZ0My5vZmZpY2VtYXRlLmNuXC9iYWNrZW5kXC9zaW5nbGVBdXRoIiwiaWF0IjoxNjUwMzQ5Mzg3LCJleHAiOjE3MjIzNDkzODcsIm5iZiI6MTY1MDM0OTM4NywianRpIjoiVGtOZk1rVjlDdDlvTkNGdCIsInN1YiI6MSwicHJ2IjoiODdlMGFmMWVmOWZkMTU4MTJmZGVjOTcxNTNhMTRlMGIwNDc1NDZhYSJ9.pZxxge_mhBRHSBeTsmZ-hrBoxEq0GKuwLnWrvjdsI9M")
//	request.Header.Set("Content-Type", helper.APPLICATION_JOSN)
//	request.Url = helper.Url{
//		Host: "http://devt3.officemate.cn",
//		Path: "backend/admin/getS4CustomerList?company_code=1000&name=东北大学",
//		Query: url.Values{
//			"company_code":{"1000"},
//			"name":{"东北大学"},
//		},
//	}
//	params := make(map[string]interface{})
//	params["company_code"] = "1020"
//	params["create_date"] = []string{"2022-04-22 00:00:00", "2022-04-22 23:59:59"}
//
//	result := request.Get()
//
//	fmt.Println(result)
//	os.Exit(1)
//}

func collyTest()  {
	c := colly.NewCollector(
		colly.AllowedDomains("learnku.com"),
	)

	c.OnHTML("div.topic-list > .simple-topic", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildAttr("div.user-avatar img", "src"))
		fmt.Println("链接：", e.ChildAttr("a.rm-tdu", "href"))
		fmt.Println(strings.Trim(strings.Replace(strings.Replace(e.ChildText("a.rm-tdu > span.topic-title"), " ", "", -1), "\n", "", -1), "new"))
	})

	//c.OnHTML("a[rel='next']", func(e *colly.HTMLElement) {
	//	fmt.Printf("%s ", e.Attr("href"))
	//	fmt.Println()
	//	e.Request.Visit(e.Attr("href"))
	//})

	c.OnHTML("a[href].page-link", func(e *colly.HTMLElement) {
		fmt.Printf("%s ", e.Attr("href"))
		fmt.Println()
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("visting", request.URL.String())
	})

	// htmlCallbacks
	c.Visit("https://learnku.com/go?order=recent&page=1")
}
