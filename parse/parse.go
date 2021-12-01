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
	s.Children().Each(func(i int, sc *goquery.Selection) {
		attr := make(map[string]string)
		if sc.Is("span") {
			// p = Piece{NORMAL_TEXT, s.Text(), nil}
			piece = append(piece, Piece{NORMAL_TEXT, sc.Text(), nil})
		} else if sc.Is("a") {
			attr["href"], _ = sc.Attr("href")
			// p = Piece{LINK, removeBrAndBlank(s.Text()), attr}
			piece = append(piece, Piece{LINK, removeBrAndBlank(sc.Text()), attr})
		} else if sc.Is("img") {
			attr["src"], _ = sc.Attr("data-src")
			attr["alt"], _ = sc.Attr("alt")
			attr["title"], _ = sc.Attr("title")
			// p = Piece{IMAGE, "", attr}
			piece = append(piece, Piece{IMAGE, "", attr})
		} else if sc.Is("ol") {
			piece = append(piece, parseList(sc, O_LIST)...)
		} else if sc.Is("ul") {
			piece = append(piece, parseList(sc, U_LIST)...)
		} else if sc.Is("section") {
			piece = append(piece, parseSection(sc)...)
		} else {
			piece = append(piece, Piece{NORMAL_TEXT, sc.Text(), nil})
		}
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
	s.Find("code").Each(func(i int, sc *goquery.Selection) {
		codeRows = append(codeRows, sc.Text())
	})
	p := Piece{CODE_BLOCK, codeRows, nil}
	return []Piece{p}
}

func parseList(s *goquery.Selection, ptype PieceType) []Piece {
	var list []Piece
	s.Find("li").Each(func(i int, sc *goquery.Selection) {
		list = append(list, Piece{ptype, parseSection(sc), nil})
	})
	return list
}

func parseBlockQuote(s *goquery.Selection) []Piece {
	var bq []Piece
	s.Children().Each(func(i int, sc *goquery.Selection) {
		bq = append(bq, Piece{BLOCK_QUOTES, parseSection(sc), nil})
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
			pieces = append(pieces, parseList(s, O_LIST)...)
		} else if s.Is("ul") {
			pieces = append(pieces, parseList(s, U_LIST)...)
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
