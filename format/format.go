package format

import (
	"strconv"

	"github.com/fengxxc/wechatmp2markdown/parse"
)

func Format(article parse.Article) string {
	var result string
	var titleMdStr string = formatTitle(article.Title)
	result += titleMdStr
	var metaMdStr string = formatMeta(article.Meta)
	result += metaMdStr
	var tagsMdStr string = formatTags(article.Tags)
	result += tagsMdStr
	var content string = formatContent(article.Content, 0)
	result += content
	return result
}

func formatTitle(piece parse.Piece) string {
	var prefix string
	level, _ := strconv.Atoi(piece.Attrs["level"])
	for i := 0; i < level; i++ {
		prefix += "#"
	}
	return prefix + " " + piece.Val.(string) + "  \n"
}

func formatMeta(meta string) string {
	return meta + "  \n" // TODO
}

func formatTags(tags string) string {
	return tags + "  \n" // TODO
}

func formatContent(pieces []parse.Piece, depth int) string {
	var contentMdStr string
	for _, piece := range pieces {
		var pieceMdStr string
		switch piece.Type {
		case parse.HEADER:
			pieceMdStr = formatTitle(piece)
		case parse.LINK:
			pieceMdStr = formatLink(piece)
		case parse.NORMAL_TEXT:
			pieceMdStr = piece.Val.(string)
		case parse.BOLD_TEXT:
			pieceMdStr = "**" + piece.Val.(string) + "**"
		case parse.ITALIC_TEXT:
			pieceMdStr = "*" + piece.Val.(string) + "*"
		case parse.BOLD_ITALIC_TEXT:
			pieceMdStr = "***" + piece.Val.(string) + "***"
		case parse.IMAGE:
			pieceMdStr = formatImage(piece)
		case parse.TABLE:
			// TODO
		case parse.CODE_INLINE:
		case parse.CODE_BLOCK:
			pieceMdStr = formatCodeBlock(piece)
		case parse.BLOCK_QUOTES:
			pieceMdStr = formatBlockQuote(piece, depth)
		case parse.O_LIST:
			pieceMdStr = formatList(piece, depth)
		case parse.U_LIST:
			pieceMdStr = formatList(piece, depth)
		case parse.HR:
		case parse.BR:
			pieceMdStr = "  \n"
		}
		contentMdStr += pieceMdStr
	}
	return contentMdStr
}

func formatBlockQuote(piece parse.Piece, depth int) string {
	var bqMdString string
	var prefix string = ">"
	for i := 0; i < depth; i++ {
		prefix += ">"
	}
	prefix += " "
	bqMdString = prefix + formatContent(piece.Val.([]parse.Piece), depth+1) + "  \n"
	return bqMdString
}

func formatList(li parse.Piece, depth int) string {
	var listMdString string
	var prefix string
	for j := 0; j < depth; j++ {
		prefix += "    "
	}
	if li.Type == parse.U_LIST {
		prefix += "- "
	} else if li.Type == parse.O_LIST {
		prefix += strconv.Itoa(1) + ". " // 写死成1也大丈夫，markdown会自动累加序号
	}
	listMdString = prefix + formatContent(li.Val.([]parse.Piece), depth+1) + "  \n"
	return listMdString
}

func formatCodeBlock(piece parse.Piece) string {
	var codeMdStr string
	codeMdStr += "```\n"
	codeRows := piece.Val.([]string)
	for _, row := range codeRows {
		codeMdStr += row + "\n"
	}
	codeMdStr += "```  \n"
	return codeMdStr
}

func formatImage(piece parse.Piece) string {
	return "![" + piece.Attrs["alt"] + "](" + piece.Attrs["src"] + " \"" + piece.Attrs["title"] + "\")"
}

func formatLink(piece parse.Piece) string {
	var linkMdStr string = "[" + piece.Val.(string) + "](" + piece.Attrs["href"] + ")"
	return linkMdStr
}
