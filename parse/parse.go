package parse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseSection(s *goquery.Selection) Paragraph {
	// fmt.Printf("s.Length() = %d\n", s.Length())
	// fmt.Printf("s.Size() = %d\n", s.Size())
	// var piece = make([]Token, s.Size())
	var piece []Piece
	s.Children().Each(func(i int, s *goquery.Selection) {
		var p Piece
		attr := make(map[string]string)
		if s.Is("span") {
			p = Piece{NORMAL_TEXT, s.Text(), nil}
		} else if s.Is("a") {
			attr["href"], _ = s.Attr("href")
			p = Piece{LINK, removeBrAndBlank(s.Text()), attr}
		} else if s.Is("img") {
			attr["src"], _ = s.Attr("src")
			attr["alt"], _ = s.Attr("alt")
			attr["title"], _ = s.Attr("title")
			p = Piece{IMAGE, "", attr}
		} else if s.Is("ol") {
			// TODO
		} else if s.Is("ul") {
			// TODO
		} else {
			p = Piece{NORMAL_TEXT, s.Text(), nil}
			// TODO
		}
		// fmt.Printf("i = %d\n", i)
		// fmt.Printf("%+v\n", t)
		// tokens[i] = t
		piece = append(piece, p)
	})
	return Paragraph{piece}
}

func parseHeader(s *goquery.Selection) Paragraph {
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
	fmt.Println("***********" + strconv.Itoa(level))
	attr := map[string]string{"level": strconv.Itoa(level)}
	p := Piece{HEADER, removeBrAndBlank(s.Text()), attr}
	return Paragraph{[]Piece{p}}
}

func parsePre(s *goquery.Selection) Paragraph {
	var codeRows []string
	s.Find("code").Each(func(i int, s *goquery.Selection) {
		codeRows = append(codeRows, s.Text())
	})
	p := Piece{CODE_BLOCK, codeRows, nil}
	return Paragraph{[]Piece{p}}
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
	fmt.Println(title)
	attr := map[string]string{"level": "1"}
	article.Title = Piece{HEADER, title, attr}

	// meta 细节待完善
	meta := mainContent.Find("#meta_content").Text()
	meta = removeBrAndBlank(meta)
	fmt.Println(meta)
	article.Meta = meta

	// tags 细节待完善
	tags := mainContent.Find("#js_tags").Text()
	tags = removeBrAndBlank(tags)
	fmt.Println(tags)
	article.Tags = tags

	// content
	// section[style="line-height: 1.5em;"]>span,a	=> 一般段落（含文本和超链接）
	// p[style="line-height: 1.5em;"]				=> 项目列表（有序/无序）
	// section[style=".*text-align:center"]>img		=> 居中段落（图片）
	content := mainContent.Find("#js_content")
	var sections []Paragraph
	content.Children().Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		// fmt.Println(s.Attr("style"))
		var paragraph Paragraph
		if s.Is("pre") || s.Is("section.code-snippet__fix") {
			// 代码块
			paragraph = parsePre(s)
		} else if s.Is("p") || s.Is("section") {
			paragraph = parseSection(s)
		} else if s.Is("h1") || s.Is("h2") || s.Is("h3") || s.Is("h4") || s.Is("h5") || s.Is("h6") {
			paragraph = parseHeader(s)
		} else if s.Is("ol") {
			// TODO
		} else if s.Is("ul") {
			// TODO
		}
		// sections[i] = block
		sections = append(sections, paragraph)
	})
	article.Content = sections

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
