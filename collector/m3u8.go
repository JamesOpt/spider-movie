package collector

import (
	"bytes"
	"fmt"
	"github.com/grafov/m3u8"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"spider-movie/app"
	"spider-movie/db"
	"spider-movie/helper"
	"spider-movie/model"
	"strconv"
	"strings"
	"sync"
)


/**
将MediaSegment内nil去除
 */
func getRealMediaPlaylist(arr []*m3u8.MediaSegment) []*m3u8.MediaSegment {
	var data []*m3u8.MediaSegment

	for _, v := range arr {
		if v == nil {
			continue
		}

		data = append(data, v)
	}

	return data
}

/**
下载ts文件
 */
func DownloadRaw(url, serialPath string, serial *model.Series, self interface{})  {
	serial.SpiderLink = url

	client := http.Client{}
	// response.Body 只能读一次，第二次读不到内容
	response, err := client.Get(url)
	if err != nil{
		panic(err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	p, listType, err := m3u8.DecodeFrom(strings.NewReader(string(data)), true)
	if err != nil {
		panic(err)
	}

	// 创建集目录
	serialAbsPath := filepath.Join(app.GetRootPath("video"), serialPath, strconv.Itoa(serial.Serial))
	os.MkdirAll(serialAbsPath, 0644)

	// 下载m3u8文件
	m3u8Filenames := strings.Split(url, "/")
	m3u8Filename := filepath.Join(serialAbsPath, m3u8Filenames[len(m3u8Filenames) - 1])

	file , _:= os.OpenFile(m3u8Filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer file.Close()

	bReader := bytes.NewReader([]byte(data))
	_, err = bReader.WriteTo(file)
	// end

	var downNum int64 = 0

	// 分析m3u8文件类型
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		segments := getRealMediaPlaylist(mediapl.Segments)
		wg := &sync.WaitGroup{}
		bar := helper.NewProgress(0, int64(len(segments)), serialPath + "【" + strconv.Itoa(serial.Serial) + "】")

		maxNum := app.Config.GetInt("coroutine_num") // 控制好携程的数量, 一次性生成多个携程有问题

		ch := make(chan int, maxNum)

		m := sync.Map{}
		for _, v := range segments {
			wg.Add(1)
			ch <- 1
			go func(segment *m3u8.MediaSegment, ch chan int) {
				var matchUrl string
				switch self.(type) {
					case Hktv,*Hktv:
						matchUrl2 := strings.Split(segment.URI, "/")
						matchUrl = matchUrl2[len(matchUrl2) -1]
					default:
						matchUrl = regexp.MustCompile(`[^\/.*]+(\.ts)`).FindString(segment.URI)
				}

				spiderm3u8 := model.SpiderM3u8{
					SeriesId: serial.ID,
					Filename: matchUrl,
				}
				db.Engine.Driver().FirstOrCreate(&spiderm3u8, spiderm3u8)

				if spiderm3u8.Status == 1 {
					downNum += 1
					bar.Play(downNum)
					<-ch
					wg.Done()
					return
				}

				spiderm3u8.Link = segment.URI
				spiderm3u8.Status = 1

				err := helper.Download(segment.URI, matchUrl, filepath.Join(serialPath, strconv.Itoa(serial.Serial)), &m)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				downNum += 1
				bar.Play(downNum)

				db.Engine.Driver().Save(spiderm3u8)

				<-ch
				wg.Done()
			}(v, ch)
		}
		db.Engine.Driver().Save(serial)
		wg.Wait()

		// 将替换的m3u8内容写入文件内
		m.Range(func(key, value any) bool {
			data = regexp.MustCompile(key.(string)).ReplaceAllFunc(data, func(i []byte) []byte {
				return []byte(value.(string))
			})

			return true
		})

		file.Seek(0, 0)
		bReader := bytes.NewReader(data)
		bReader.WriteTo(file)

	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("1111 %+v\n", masterpl.String())
	}
}