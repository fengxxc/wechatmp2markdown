package parse

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseSection(s *goquery.Selection) []Piece {
	var piece []Piece
	s.Children().Each(func(i int, s *goquery.Selection) {
		var p []Piece
		attr := make(map[string]string)
		if s.Is("span") {
			// p = Piece{NORMAL_TEXT, s.Text(), nil}
			p = append(p, Piece{NORMAL_TEXT, s.Text(), nil})
		} else if s.Is("a") {
			attr["href"], _ = s.Attr("href")
			// p = Piece{LINK, removeBrAndBlank(s.Text()), attr}
			p = append(p, Piece{LINK, removeBrAndBlank(s.Text()), attr})
		} else if s.Is("img") {
			attr["src"], _ = s.Attr("data-src")
			attr["alt"], _ = s.Attr("alt")
			attr["title"], _ = s.Attr("title")
			// p = Piece{IMAGE, "", attr}
			p = append(p, Piece{IMAGE, "", attr})
		} else if s.Is("ol") {
			p = append(p, parseOl(s)...)
		} else if s.Is("ul") {
			p = append(p, parseUl(s)...)
		} else if s.Is("section") {
			p = append(p, parseSection(s)...)
		} else {
			p = append(p, Piece{NORMAL_TEXT, s.Text(), nil})
		}
		piece = append(piece, p...)
	})
	return piece
}

func parseHeader(s *goquery.Selection) []Piece {
	var level int
	switch {
	case s.Is("h1"):
		level = 1
	case s.Is("h2"):
		level = 2
	case s.Is("h3"):
		level = 3
	case s.Is("h4"):
		level = 4
	case s.Is("h5"):
		level = 5
	case s.Is("h6"):
		level = 6
	}
	attr := map[string]string{"level": strconv.Itoa(level)}
	p := Piece{HEADER, removeBrAndBlank(s.Text()), attr}
	return []Piece{p}
}

func parsePre(s *goquery.Selection) []Piece {
	var codeRows []string
	s.Find("code").Each(func(i int, s *goquery.Selection) {
		codeRows = append(codeRows, s.Text())
	})
	p := Piece{CODE_BLOCK, codeRows, nil}
	return []Piece{p}
}

func parseUl(s *goquery.Selection) []Piece {
	var list []Piece
	s.Find("li").Each(func(i int, s *goquery.Selection) {
		list = append(list, Piece{U_LIST, parseSection(s), nil})
	})
	return list
}

func parseOl(s *goquery.Selection) []Piece {
	var list []Piece
	s.Find("li").Each(func(i int, s *goquery.Selection) {
		list = append(list, Piece{O_LIST, parseSection(s), nil})
	})
	return list
}

func parseBlockQuote(s *goquery.Selection) []Piece {
	var bq []Piece
	s.Children().Each(func(i int, s *goquery.Selection) {
		bq = append(bq, Piece{BLOCK_QUOTES, parseSection(s), nil})
	})
	return bq
}

func ParseFromReader(r io.Reader) Article {
	var article Article
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	var mainContent *goquery.Selection = doc.Find("#img-content")

	// 标题
	title := mainContent.Find("#activity-name").Text()
	attr := map[string]string{"level": "1"}
	article.Title = Piece{HEADER, removeBrAndBlank(title), attr}

	// meta 细节待完善
	meta := mainContent.Find("#meta_content").Text()
	meta = removeBrAndBlank(meta)
	article.Meta = meta

	// tags 细节待完善
	tags := mainContent.Find("#js_tags").Text()
	tags = removeBrAndBlank(tags)
	article.Tags = tags

	// content
	// section[style="line-height: 1.5em;"]>span,a	=> 一般段落（含文本和超链接）
	// p[style="line-height: 1.5em;"]				=> 项目列表（有序/无序）
	// section[style=".*text-align:center"]>img		=> 居中段落（图片）
	content := mainContent.Find("#js_content")
	// var sections []Paragraph
	var pieces []Piece
	content.Children().Each(func(i int, s *goquery.Selection) {
		// var paragraph Paragraph
		if s.Is("pre") || s.Is("section.code-snippet__fix") {
			// 代码块
			pieces = append(pieces, parsePre(s)...)
		} else if s.Is("p") || s.Is("section") {
			pieces = append(pieces, parseSection(s)...)
		} else if s.Is("h1") || s.Is("h2") || s.Is("h3") || s.Is("h4") || s.Is("h5") || s.Is("h6") {
			pieces = append(pieces, parseHeader(s)...)
		} else if s.Is("ol") {
			pieces = append(pieces, parseOl(s)...)
		} else if s.Is("ul") {
			pieces = append(pieces, parseUl(s)...)
		} else if s.Is("blockquote") {
			pieces = append(pieces, parseBlockQuote(s)...)
		}
		// sections = append(sections, paragraph)
		pieces = append(pieces, Piece{BR, nil, nil})
	})
	// article.Content = sections
	article.Content = pieces

	return article
}

func ParseFromHTMLString(s string) Article {
	return ParseFromReader(strings.NewReader(s))
}

func ParseFromHTMLFile(filepath string) Article {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic(err)
	}
	return ParseFromReader(bytes.NewReader(content))
}

func ParseFromURL(url string) Article {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("get from url %s error: %d %s", url, res.StatusCode, res.Status)
	}
	return ParseFromReader(res.Body)
}

func removeBrAndBlank(s string) string {
	regstr := "\\s{2,}"
	reg, _ := regexp.Compile(regstr)
	sb := make([]byte, len(s))
	copy(sb, s)
	spc_index := reg.FindStringIndex(string(sb)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		sb = append(sb[:spc_index[0]+1], sb[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(sb))            //继续在字符串中搜索
	}
	return strings.Replace(string(sb), "\n", " ", -1)
}
