package parse

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseSection(s *goquery.Selection, imagePolicy ImagePolicy, lastPieceType PieceType) []Piece {
	var pieces []Piece
	if lastPieceType == O_LIST || lastPieceType == U_LIST || lastPieceType == NULL || lastPieceType == BLOCK_QUOTES {
		// pieces = append(pieces, Piece{NULL, nil, nil})
	} else {
		pieces = append(pieces, Piece{BR, nil, nil})
	}
	var _lastPieceType PieceType = NULL
	s.Contents().Each(func(i int, sc *goquery.Selection) {
		attr := make(map[string]string)
		if sc.Is("a") {
			attr["href"], _ = sc.Attr("href")
			pieces = append(pieces, Piece{LINK, removeBrAndBlank(sc.Text()), attr})
		} else if sc.Is("figure") {
			pieces = append(pieces, parseSection(sc, imagePolicy, _lastPieceType)...)
		} else if sc.Is("img") {
			attr["src"], _ = sc.Attr("data-src")
			attr["alt"], _ = sc.Attr("alt")
			attr["title"], _ = sc.Attr("title")
			switch imagePolicy {
			case IMAGE_POLICY_URL:
				pieces = append(pieces, Piece{IMAGE, nil, attr})
			case IMAGE_POLICY_SAVE:
				image := fetchImgFile(attr["src"])
				pieces = append(pieces, Piece{IMAGE, image, attr})
			case IMAGE_POLICY_BASE64:
				fallthrough
			default:
				base64Image := img2base64(fetchImgFile(attr["src"]))
				pieces = append(pieces, Piece{IMAGE_BASE64, base64Image, attr})
			}
		} else if sc.Is("ol") {
			pieces = append(pieces, parseList(sc, O_LIST, imagePolicy)...)
		} else if sc.Is("ul") {
			pieces = append(pieces, parseList(sc, U_LIST, imagePolicy)...)
		} else if sc.Is("pre") || sc.Is("section.code-snippet__fix") {
			// 代码块
			pieces = append(pieces, parsePre(sc)...)
		} else if sc.Is("span") {
			pieces = append(pieces, parseSection(sc, imagePolicy, _lastPieceType)...)
		} else if sc.Is("p") || sc.Is("section") {
			pieces = append(pieces, parseSection(sc, imagePolicy, _lastPieceType)...)
			if removeBrAndBlank(sc.Text()) != "" && len(pieces) > 0 && pieces[len(pieces)-1].Type != BR {
				pieces = append(pieces, Piece{BR, nil, nil})
			}
		} else if sc.Is("h1") || sc.Is("h2") || sc.Is("h3") || sc.Is("h4") || sc.Is("h5") || sc.Is("h6") {
			pieces = append(pieces, parseHeader(sc)...)
		} else if sc.Is("blockquote") {
			pieces = append(pieces, parseBlockQuote(sc, imagePolicy)...)
		} else if sc.Is("strong") {
			pieces = append(pieces, parseStrong(sc)...)
		} else if sc.Is("table") {
			pieces = append(pieces, parseTable(sc)...)
		} else {
			if sc.Text() != "" {
				pieces = append(pieces, Piece{NORMAL_TEXT, sc.Text(), nil})
			}
		}
		if len(pieces) > 0 {
			_lastPieceType = pieces[len(pieces)-1].Type
		}
	})
	return pieces
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
	// TODO when include img...
	var codeRows []string
	s.Find("code").Each(func(i int, sc *goquery.Selection) {
		var codeLine string = ""
		sc.Contents().Each(func(i int, sc *goquery.Selection) {
			if goquery.NodeName(sc) == "br" {
				codeRows = append(codeRows, codeLine)
				codeLine = ""
			} else {
				codeLine += sc.Text()
			}
		})
		codeRows = append(codeRows, codeLine)
	})
	p := Piece{CODE_BLOCK, codeRows, nil}
	return []Piece{p}
}

func parseList(s *goquery.Selection, ptype PieceType, imagePolicy ImagePolicy) []Piece {
	var list []Piece
	s.Find("li").Each(func(i int, sc *goquery.Selection) {
		list = append(list, Piece{ptype, parseSection(sc, imagePolicy, ptype), nil})
	})
	return list
}

func parseBlockQuote(s *goquery.Selection, imagePolicy ImagePolicy) []Piece {
	var bq []Piece
	s.Contents().Each(func(i int, sc *goquery.Selection) {
		bq = append(bq, Piece{BLOCK_QUOTES, parseSection(sc, imagePolicy, BLOCK_QUOTES), nil})
	})
	bq = append(bq, Piece{BR, nil, nil})
	return bq
}

func parseTable(s *goquery.Selection) []Piece {
	// 先简单粗暴把原生的挪过去
	var table []Piece
	html, _ := s.Html()
	table = append(table, Piece{TABLE, "<table>" + html + "</table>", map[string]string{"type": "native"}})
	return table
}

func parseStrong(s *goquery.Selection) []Piece {
	var bt []Piece
	bt = append(bt, Piece{BOLD_TEXT, strings.TrimSpace(s.Text()), nil})
	return bt
}

func parseMeta(s *goquery.Selection) []string {
	var res []string
	s.Children().Each(func(i int, sc *goquery.Selection) {
		if sc.Is("#profileBt") {
			res = append(res, removeBrAndBlank(sc.Find("#js_name").Text()))
		} else {
			style, exists := sc.Attr("style")
			if !(exists && strings.Contains(style, "display: none;")) {
				// t := sc.Nodes[0].Data
				t := strings.TrimSpace(sc.Text())
				res = append(res, t)
			}
		}
	})
	return res
}

func ParseFromReader(r io.Reader, imagePolicy ImagePolicy) Article {
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

	// meta
	meta := mainContent.Find("#meta_content")
	metastring := parseMeta(meta)
	article.Meta = metastring
	// 从js中找到发布时间
	re, _ := regexp.Compile("var ct = \"([0-9]+)\"")
	findstrs := re.FindStringSubmatch(doc.Find("script").Text())
	if len(findstrs) > 1 {
		var createTime string = findstrs[1]
		timestamp, _ := strconv.Atoi(createTime)
		time := time.Unix(int64(timestamp), 0)
		article.Meta = append(article.Meta, time.Format("2006-01-02 15:04"))
	}

	// tags 细节待完善
	tags := mainContent.Find("#js_tags").Text()
	tags = removeBrAndBlank(tags)
	article.Tags = tags

	// content
	// section[style="line-height: 1.5em;"]>span,a	=> 一般段落（含文本和超链接）
	// p[style="line-height: 1.5em;"]				=> 项目列表（有序/无序）
	// section[style=".*text-align:center"]>img		=> 居中段落（图片）
	content := mainContent.Find("#js_content")
	pieces := parseSection(content, imagePolicy, NULL)
	article.Content = pieces

	return article
}

func ParseFromHTMLString(s string, imagePolicy ImagePolicy) Article {
	return ParseFromReader(strings.NewReader(s), imagePolicy)
}

func ParseFromHTMLFile(filepath string, imagePolicy ImagePolicy) Article {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err2 := io.ReadAll(file)
	if err2 != nil {
		panic(err)
	}
	return ParseFromReader(bytes.NewReader(content), imagePolicy)
}

func ParseFromURL(url string, imagePolicy ImagePolicy) Article {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("get from url %s error: %d %s", url, res.StatusCode, res.Status)
	}
	return ParseFromReader(res.Body, imagePolicy)
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

func fetchImgFile(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("get Image from url %s error: %s", url, err.Error())
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("get Image from url %s error: %d %s", url, res.StatusCode, res.Status)
	}
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("read image Response error: %s", err.Error())
	}
	return content
}

func img2base64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}

type ImagePolicy int32

const (
	IMAGE_POLICY_URL ImagePolicy = iota
	IMAGE_POLICY_SAVE
	IMAGE_POLICY_BASE64
)

func ImageArgValue2ImagePolicy(val string) ImagePolicy {
	var imagePolicy ImagePolicy
	switch val {
	case "url":
		imagePolicy = IMAGE_POLICY_URL
	case "save":
		imagePolicy = IMAGE_POLICY_SAVE
	case "base64":
		fallthrough
	default:
		imagePolicy = IMAGE_POLICY_BASE64
	}
	return imagePolicy
}
