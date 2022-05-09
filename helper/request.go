package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"runtime"
	"spider-movie/app"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	APPLICATION_JOSN = "application/json"
)

type Url struct {
	Host string
	Path string
	Query netUrl.Values
}

type Request struct {
	//Url
	Header http.Header
	Client *http.Client
}

func (url *Url) check() error {
	if url.Host == "" || url.Path == "" {
		return errors.New("Host or Path is not null")
	}

	return nil
}

func (url *Url) makeUrl() string {
	path := ""
	if 0 != strings.Index("/", url.Path) {
		path = "/" + url.Path
	} else {
		path = url.Path
	}

	if url.Query != nil {
		path += "?" + url.Query.Encode()
	}

	return url.Host + path
}

func NewRequest() *Request {
	return &Request{
		Header: make(http.Header),
		Client: nil,
	}
}

func (req *Request) Do(method string, api string, body io.Reader) io.ReadCloser {
	//if err := req.Url.check(); err != nil {
	//	panic(err)
	//}

	newReq, err := http.NewRequest(method, api, body)

	if err != nil {
		panic(err)
	}

	newReq.Header = req.Header

	if req.Client == nil {
		req.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 0 && via[0].Proto != "HTTP/2.0"{
					c := http.Client{}
					newReq, _ := http.NewRequest("HEAD", req.URL.String(), nil)
					response, _ := c.Do(newReq)
					defer response.Body.Close()

					if response.Proto == "HTTP/2.0" {
						return errors.New("HTTP/2.0 request cache error")
					}
				}

				return nil
			},
		}
	}

	response, err := req.Client.Do(newReq)

	if err != nil {
		if err.(*netUrl.Error).Err.Error() != errors.New("HTTP/2.0 request cache error").Error() {
			panic(err)
		}
	}

	if response.Header.Get("Location") != ""{
		return NewRequest().Do(method, response.Header.Get("Location"), body)
	}

	return response.Body
}

func (req *Request) Get(api string) string {

	body := req.Do("GET", api, nil)
	defer body.Close()

	dom, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(dom)
}

func (req *Request) Post(api string, params map[string]interface{}) string {
	data, err := json.Marshal(params)

	if err != err{
		panic(err)
	}

	body := req.Do("POST", api, strings.NewReader(string(data)))
	defer body.Close()

	dom, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(dom)
}

func Download(uri string, filename string, basePath interface{}, m3u8Filename string) error {
	//u, err := netUrl.Parse(uri)
	//if err != nil {
	//	return err
	//}


	//originFilename := filename

	serialPath := app.GetRootPath("video")
	absPath := filepath.Join(serialPath, basePath.(string))

	os.MkdirAll(absPath, 0644)

	// 如果文件名字超过255，则换名
	if utf8.RuneCountInString(filename) > 255 {
		fmt.Println(filename)
		filename = randStrings(12)
		fmt.Println(filename)

		//m3File, _ := os.OpenFile(m3u8Filename, os.O_RDWR|os.O_EXCL, 0644)
		//
		//defer m3File.Close()
		//m3Data, _ := io.ReadAll(m3File)
		//
		//newM3Data := regexp.MustCompile(originFilename).ReplaceAllFunc(m3Data, func(i []byte) []byte {
		//	fmt.Println(filename)
		//	return []byte(tfilename)
		//})
		//
		//m3File.Write(newM3Data)
	}

	file , err:= os.OpenFile(filepath.Join(absPath, filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	// 设置参数
	response := NewRequest().Get(uri)

	reader := bytes.NewReader([]byte(response))
	_, err = reader.WriteTo(file)

	if err != nil{
		panic(err)
	}

	return nil
}

//func (req *Request) SetUrl(host, path string, query netUrl.Values) *Request {
//	req.Url = Url{
//		Host:  host,
//		Path:  path,
//		Query: query,
//	}
//
//	return req
//}

/**
随机生成字符串
 */
func randStrings(length  int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id := curGoroutineId()

	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length - len(id))
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result) + id
}

/**
获取当前协程id
 */
func curGoroutineId() string {
	buf := make([]byte, 64)
	runtime.Stack(buf, false)

	b := bytes.TrimPrefix(buf, []byte("goroutine "))
	index := bytes.IndexByte(b, ' ')
	if index < 0 {
		return "0"
	} else {
		b = b[:index]
	}

	return string(b)
}