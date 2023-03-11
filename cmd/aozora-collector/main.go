package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 作品情報をまとめるstruct
type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	InfoURL  string
	ZipURL   string
}

// goqueryはスクレイピングできるライブラリ
// goqueryでURLからDOMオブジェクトを作成する
func findEntries(siteURL string) ([]Entry, error) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	// リンクURL一覧を取得する
	pat := regexp.MustCompile(`.*/cards/([0-9]+/card([0-9]+).html$)`)
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/cards%s.html", token[1], token[2])
		// 作者とZIPファイルのURLを得る
		author, zipURL := findAuthorAndZIP(pageURL)
		println(author, zipURL)
	})
	return nil, nil
}

// 作者とZIPファイルを得る関数の雛形
func findAuthorAndZIP(siteURL string) (string, string) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return "", ""
	}

	author := doc.Find("table[sammary=作家データ] tr:nth-child(1) td:nth-child(2)").Text()
	zipURL := ""
	doc.Find("table.downloada").Each(func(n int, elem *goquery.Selection) {
		href := elem.AttrOr("href", "")
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})
	return author, zipURL
}

// zipファイルのURL一覧を取得する
func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Title, entry.ZipURL)
	}
}
