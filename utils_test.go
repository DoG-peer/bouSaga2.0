package main

import (
	"testing"
)

func TestSplitToRes(t *testing.T) {
	data := `
1<>名無しさん<>sage<>now<>hogehoge<><>
	`
	value := SplitToRes(data)
	expected := []Res{
		Res{id: 1, name: "名無しさん", email: "sage", date: "now", body: "hogehoge"},
	}
	if value[0] != expected[0] {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}

func TestMaxID(t *testing.T) {
	data := []Res{
		Res{id: 1, name: "名無しさん", email: "sage", date: "now", body: "hogehoge"},
		Res{id: 2, name: "名無しさん", email: "sage", date: "now", body: "hogehoge"},
		Res{id: 100, name: "名無しさん", email: "sage", date: "now", body: "hogehoge"},
	}
	value := MaxID(data, 0)
	expected := 100
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}

func TestGetRawURL(t *testing.T) {
	url := "hoge/123/456/"
	value, e := GetRawURL(url)
	expected := "http://jbbs.shitaraba.net/bbs/rawmode.cgi/game/123/456/"
	if e != nil {
		t.Fatalf("Failed to parse %v", e)
	} else if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}
