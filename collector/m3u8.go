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
	fmt.Println(url)
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

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		segments := getRealMediaPlaylist(mediapl.Segments)
		wg := &sync.WaitGroup{}
		for _, v := range segments {
			wg.Add(1)
			go func(segment *m3u8.MediaSegment) {
				matchUrl := regexp.MustCompile(`[^\/.*]+\.ts`).FindString(segment.URI)
				helper.Download(v.URI, matchUrl, filepath.Join(basePath, strconv.Itoa(serial.Serial)))
				wg.Done()
			}(v)
		}
		db.Engine.Driver().Save(serial)
		wg.Wait()

	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("1111 %+v\n", masterpl.String())
	}
}