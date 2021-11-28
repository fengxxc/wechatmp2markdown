package parse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseSection(s *goquery.Selection) Block {
	// fmt.Printf("s.Length() = %d\n", s.Length())
	// fmt.Printf("s.Size() = %d\n", s.Size())
	// var tokens = make([]Token, s.Size())
	var tokens []Token
	s.Children().Each(func(i int, s *goquery.Selection) {
		var t Token
		attr := make(map[string]string)
		if s.Is("span") {
			t = Token{NORMAL_TEXT, s.Text(), nil}
		} else if s.Is("a") {
			attr["href"], _ = s.Attr("href")
			t = Token{LINK, removeBrAndBlank(s.Text()), attr}
		} else if s.Is("img") {
			attr["src"], _ = s.Attr("src")
			t = Token{IMAGE, "", attr}
		} else {
			t = Token{NORMAL_TEXT, s.Text(), nil}
			// TODO
		}
		// fmt.Printf("i = %d\n", i)
		// fmt.Printf("%+v\n", t)
		// tokens[i] = t
		tokens = append(tokens, t)
	})
	return Block{tokens}
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
	article.title = title

	// meta 细节待完善
	meta := mainContent.Find("#meta_content").Text()
	meta = removeBrAndBlank(meta)
	fmt.Println(meta)
	article.meta = meta

	// tags 细节待完善
	tags := mainContent.Find("#js_tags").Text()
	tags = removeBrAndBlank(tags)
	fmt.Println(tags)
	article.tags = tags

	// content
	// section[style="line-height: 1.5em;"]>span,a	=> 一般段落（含文本和超链接）
	// p[style="line-height: 1.5em;"]				=> 项目列表（有序/无序）
	// section[style=".*text-align:center"]>img		=> 居中段落（图片）
	content := mainContent.Find("#js_content")
	var sections []Block
	content.Find("section,p").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		// fmt.Println(s.Attr("style"))
		var block Block
		if s.Is("p") {
			block = parseSection(s)
		} else if s.Is("section") {
			block = parseSection(s)
		} else {
			// TODO
		}
		// sections[i] = block
		sections = append(sections, block)
	})
	article.content = sections

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
