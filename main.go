package main

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/gocolly/colly"
	"net/url"
	"os"
	"spider-movie/collector"
	"spider-movie/helper"
	"strings"
)

type Hash []byte

type Magnet struct {
	InfoHash Hash
	Trackers []string
	DisplayName string
}

const xtPrefix  = "urn:btih:"

func (m Magnet) String() string {
	ret := "magnet:?xt="
	ret += xtPrefix + hex.EncodeToString(m.InfoHash[:])
	if m.DisplayName != "" {
		ret += "&dn=" + url.QueryEscape(m.DisplayName)
	}
	for _, tr := range m.Trackers {
		ret += "&tr=" + url.QueryEscape(tr)
	}
	return ret
}

func ParseMagnetURI(uri string) (m Magnet, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		err = fmt.Errorf("error parsing uri : %s", err)
		return
	}
	if u.Scheme != "magnet" {
		err = fmt.Errorf("unexpected scheme: %q", u.Scheme)
		return
	}
	xt := u.Query().Get("xt")
	if !strings.HasPrefix(xt, xtPrefix) {
		err = fmt.Errorf("bad xt parameter")
		return
	}


	infoHash := xt[len(xtPrefix):]


	var decode func(dst, src []byte) (int, error)
	switch len(infoHash) {
	case 40:
		decode = hex.Decode
	case 32:
		decode = base32.StdEncoding.Decode
	}

	if decode == nil {
		err = fmt.Errorf("unhandled xt parameter encoding: encoded lenght", len(infoHash))
		return
	}

	m.InfoHash = make([]byte, len(infoHash))
	n, err := decode(m.InfoHash[:], []byte(infoHash))
	if err != nil {
		err = fmt.Errorf("error decoding xt: %s", err)
		return
	}
	fmt.Println(n)
	os.Exit(1)

	if n != 20 {
		panic(n)
	}
	m.DisplayName = u.Query().Get("dn")
	m.Trackers = u.Query()["tr"]
	return
}

func main(){
	req := helper.NewRequest()
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36 Edg/101.0.1210.32")
	req.Header.Set("authority", "jx2api.5408h.cn")

	r := req.Get("https://jx2api.5408h.cn/hls/DD614ADD7D79BF27FD2B493D74B385A4C4DBBAEE48BA4A9E3D2847E6DB6BCB01880F2977A59D9238986528EF4F8D5592311E702F66228CCC36E98F9517451D34C9967201ECDE100230AE592B4F96062A5BCF6D5C2759431A7F162EED83C4C747B48F9CD3C1ABAAB68EFD9C8AF3444B85D21D91597CEBA026A319B8ABFD03413763A438DE64DC022D309705F81D13")
	//r := req.Get("https://jx2api.5408h.cn/hls/DD614ADD7D79BF27F738546065A8CAABC9D6BAF451B25D9029244EA1DF6BD3008D086A75F888923F897026EF14BB519E371C3A39075AD9A769B5CCC140775D578ED60555D9A16C0E00D93F6949D80A0A0BB3112365731B457E3103")

	fmt.Println(r)
	os.Exit(1)
	//file, _ := os.OpenFile("1", os.O_RDWR, 0644)
	//
	//defer file.Close()
	//data, _ := ioutil.ReadAll(file)
	//
	//file.Seek(0, 0)
	//file.Write([]byte{255, 255, 255,255})
	//
	//for i := 0; i < 10; i++ {
	//	if i % 2 == 0 {
	//		fmt.Println()
	//	}
	//	fmt.Printf("%s", hex.EncodeToString([]byte{data[i]}))
	//}
	//
	//
	//os.Exit(1)
	//bytes.NewReader()


// https://0ranga.com/2018/08/26/bt-metadata/

	//magn := "magnet:?xt=urn:btih:1a84227232a032c872a5e4e1432d72d167c57544&dn=[%E7%94%B5%E5%BD%B1%E5%A4%A9%E5%A0%82www.dytt89.com]%E6%96%B0%E8%9D%99%E8%9D%A0%E4%BE%A0-2022_HD%E4%B8%AD%E8%8B%B1%E5%8F%8C%E5%AD%97.mp4"
	//
	//mag, _ := ParseMagnetURI(magn)
	//fmt.Println(mag.DisplayName, mag.Trackers, string(mag.InfoHash))

	//c := collector.PianBa{
	//	Request: helper.NewRequest(),
	//}
	//c.Run("https://www.pianba.tv/html/209094.html")

	c := collector.Hktv{
		Request: helper.NewRequest(),
	}

	c.Run("https://www.hktv03.com/vod/detail/id/182142.html")
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
