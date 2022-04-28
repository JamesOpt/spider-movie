package collector

import (
	"fmt"
	"github.com/grafov/m3u8"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
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

func DownloadRaw(url, basePath string, serial *model.Series)  {
	serial.SpiderLink = url

	client := http.Client{}
	response, err := client.Get(url)
	if err != nil{
		panic(err)
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	reader := strings.NewReader(string(data))
	p, listType, err := m3u8.DecodeFrom(reader, true)
	if err != nil {
		panic(err)
	}

	var downNum int64 = 0

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		segments := getRealMediaPlaylist(mediapl.Segments)
		wg := &sync.WaitGroup{}
		bar := helper.NewProgress(0, int64(len(segments)), basePath + "【" + strconv.Itoa(serial.Serial) + "】")

		maxNum := 50
		ch := make(chan int, maxNum)

		for _, v := range segments {
			wg.Add(1)
			ch <- 1
			go func(segment *m3u8.MediaSegment, ch chan int) {
				matchUrl := regexp.MustCompile(`[^\/.*]+\.ts`).FindString(segment.URI)
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

				err := helper.Download(segment.URI, matchUrl, filepath.Join(basePath, strconv.Itoa(serial.Serial)))
				if err != nil {
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

	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("1111 %+v\n", masterpl.String())
	}
}