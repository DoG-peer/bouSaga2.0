package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// EucToUtf8 convert euc encoded string to utf-8 encoded string
func EucToUtf8(data string) (string, error) {
	in := bytes.NewBufferString(data)
	out := new(bytes.Buffer)
	reader := transform.NewReader(in, japanese.EUCJP.NewDecoder())
	_, e := io.Copy(out, reader)
	return out.String(), e
}

// Res is
type Res struct {
	id    int
	name  string
	email string
	date  string
	body  string
}

// SplitToRes split data string to a slice of Res
func SplitToRes(data string) []Res {
	rs := []Res{}
	for _, r := range strings.Split(data, "\n") {
		ar := strings.Split(r, "<>")
		if len(ar) == 7 {
			id, _ := strconv.Atoi(ar[0])
			rs = append(rs, Res{
				id:    id,
				name:  ar[1],
				email: ar[2],
				date:  ar[3],
				body:  ar[4],
			})
		}
	}
	return rs
}

// ReadBBS gets responses of the url
func ReadBBS(url string) ([]Res, error) {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	arr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data, err := EucToUtf8(string(arr))
	if err != nil {
		return nil, err
	}

	ls := SplitToRes(data)
	return ls, nil
}

// MaxID returns maximum id or default id
func MaxID(rs []Res, k int) int {
	n := k
	for _, r := range rs {
		if n < r.id {
			n = r.id
		}
	}
	return n
}

// NextID returns maximum id or default id
func NextID(rs []Res, k int) int {
	n := k - 1
	for _, r := range rs {
		if n < r.id {
			n = r.id
		}
	}
	return n + 1
}

// GetRawURL convert usual thread url to raw url
func GetRawURL(url string) (string, error) {
	r := regexp.MustCompile(`^.*/(\d+)/(\d+)/$`)
	if !r.MatchString(url) {
		return "", fmt.Errorf("%s should be like ***/1234/5678/", url)
	}
	rawmodeURLFormat := "http://jbbs.shitaraba.net/bbs/rawmode.cgi/game/$1/$2/"
	return r.ReplaceAllString(url, rawmodeURLFormat), nil
}
