package main

import (
	"strconv"
)

// BBS communicates with Shitaraba
type BBS struct {
	baseurl string
	apiurl  string
	num     int
}

func (bbs *BBS) Read() ([]Res, error) {
	return ReadBBS(bbs.apiurl)
}

// MoveTo new url
func (bbs *BBS) MoveTo(n int) {
	bbs.apiurl = bbs.baseurl + strconv.Itoa(n) + "-"
	bbs.num = n
}
