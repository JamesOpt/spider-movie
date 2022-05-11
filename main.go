package main

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
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

/**
ffmpeg -allowed_extensions ALL -i index.m3u8 -c copy out.mp4
合成ts文件
 */
func main(){
	//c := collector.Hktv{
	//	Request: helper.NewRequest(),
	//}
	//
	//c.Run("https://www.hktv03.com/vod/detail/id/182142.html")

	c1 := collector.PianBa{
		Request: helper.NewRequest(),
	}

	c1.Run("https://www.pianba.tv/html/204591.html")
}

