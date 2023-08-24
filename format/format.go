package format

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/fengxxc/wechatmp2markdown/parse"
	"github.com/fengxxc/wechatmp2markdown/util"
)

// Format format article
func Format(article parse.Article) (string, map[string][]byte) {
	var result string
	var titleMdStr string = formatTitle(article.Title)
	result += titleMdStr
	var metaMdStr string = formatMeta(article.Meta)
	result += metaMdStr
	var tagsMdStr string = formatTags(article.Tags)
	result += tagsMdStr
	var saveImageBytes map[string][]byte
	content, saveImageBytes := formatContent(article.Content, 0)
	result += content
	return result, saveImageBytes
}

// windows下, 文件名包含非法字符时, 用相似的Unicode字符进行替换; 长度超过255个字符时，保留前255个字符
func legalizationFilenameForWindows(name string) string {
	// Windows文件名不能包含这些字符
	invalidChars := regexp.MustCompile(`[\\/:*?\"<>|]`)

	// 如果包含非法字符,则替换
	if invalidChars.MatchString(name) {
		name = strings.ReplaceAll(name, "<", "≺")
		name = strings.ReplaceAll(name, ">", "≻")
		name = strings.ReplaceAll(name, ":", ":")
		name = strings.ReplaceAll(name, "\"", "“")
		name = strings.ReplaceAll(name, "/", "∕")
		name = strings.ReplaceAll(name, "\\", "∖")
		name = strings.ReplaceAll(name, "|", "∣")
		name = strings.ReplaceAll(name, "?", "?")
		name = strings.ReplaceAll(name, "*", "⁎")
	}

	// 文件名最大长度255字符
	if len(name) > 255 {
		// 超出就截断文件名
		name = name[:255]
	}

	return name
}

// FormatAndSave fomat article and save to local file
func FormatAndSave(article parse.Article, filePath string) error {
	// basrPath := filepath.Join(filePath, )
	var basePath string
	var fileName string
	var isWin bool = runtime.GOOS == "windows"
	var separator string
	if isWin {
		separator = "\\"
	} else {
		separator = "/"
	}
	if filePath == "" {
		filePath = "." + separator
	}
	if strings.HasPrefix(filePath, "./") || strings.HasPrefix(filePath, ".\\") {
		wd, _ := os.Getwd()
		filePath = strings.Replace(filePath, ".", wd, 1)
	}
	if strings.HasSuffix(filePath, ".md") {
		// basePath = filePath[:len(filePath)-len(".md")]
		basePath = filePath[:strings.LastIndex(filePath, separator)]
		fileName = filePath
	} else {
		title := strings.TrimSpace(article.Title.Val.(string))
		if isWin {
			title = legalizationFilenameForWindows(title)
		}
		// title := "thisistitle"
		basePath = filepath.Join(filePath, title)
		fileName = filepath.Join(basePath, title+".md")
	}

	// make basePath dir if not exists
	if _, err := os.Stat(basePath); err != nil {
		if err := os.MkdirAll(basePath, 0644); err != nil {
			panic(err)
		}
	}

	var saveImageBytes map[string][]byte
	result, saveImageBytes := Format(article)
	if len(saveImageBytes) > 0 {
		for imgTitle := range saveImageBytes {
			// save to local
			imgfileName := filepath.Join(basePath, imgTitle)
			/* if err := ioutil.WriteFile(imgfileName, saveImageBytes[imgTitle], 0644); err != nil {
				log.Fatalf("can not save image file: %s\n err: %v", imgfileName, err)
				continue
			} */
			f, err := os.Create(imgfileName)
			if err != nil {
				// log.Fatalf("can not save image file: %s", imgTitle)
				log.Fatalf("can not save image file: %s\n err: %v", imgfileName, err)
				continue
			}
			defer f.Close()
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, saveImageBytes[imgTitle])
			f.Write(buf.Bytes())
		}
	}
	return os.WriteFile(fileName, []byte(result), 0644)
}

func formatTitle(piece parse.Piece) string {
	var prefix string
	level, _ := strconv.Atoi(piece.Attrs["level"])
	for i := 0; i < level; i++ {
		prefix += "#"
	}
	return prefix + " " + piece.Val.(string) + "  \n"
}

func formatMeta(meta []string) string {
	return strings.Join(meta, " ") + "  \n" // TODO
}

func formatTags(tags string) string {
	return tags + "  \n" // TODO
}

func formatContent(pieces []parse.Piece, depth int) (string, map[string][]byte) {
	var contentMdStr string
	var base64Imgs []string
	var saveImageBytes map[string][]byte = make(map[string][]byte)
	for _, piece := range pieces {
		var pieceMdStr string
		var patchSaveImageBytes map[string][]byte
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
			if piece.Val == nil {
				pieceMdStr = formatImageInline(piece)
			} else {
				// will save to local
				src := piece.Attrs["src"]
				imgExt := util.ParseImageExtFromSrc(src)
				var hashName string = util.MD5(piece.Val.([]byte)) + "." + imgExt
				saveImageBytes[hashName] = piece.Val.([]byte)
				pieceMdStr = formatImageFileReferInline(piece.Attrs["alt"], hashName)
			}
		case parse.IMAGE_BASE64:
			pieceMdStr = formatImageRefer(piece, len(base64Imgs))
			base64Imgs = append(base64Imgs, piece.Val.(string))
		case parse.TABLE:
			// TODO
		case parse.CODE_INLINE:
			// TODO
		case parse.CODE_BLOCK:
			pieceMdStr = formatCodeBlock(piece)
		case parse.BLOCK_QUOTES:
			pieceMdStr, patchSaveImageBytes = formatBlockQuote(piece, depth)
		case parse.O_LIST:
			pieceMdStr, patchSaveImageBytes = formatList(piece, depth)
		case parse.U_LIST:
			pieceMdStr, patchSaveImageBytes = formatList(piece, depth)
		case parse.HR:
			// TODO
		case parse.BR:
			pieceMdStr = "  \n"
		}
		contentMdStr += pieceMdStr
		util.MergeMap(saveImageBytes, patchSaveImageBytes)
	}
	for i := 0; i < len(base64Imgs); i++ {
		contentMdStr += "\n[" + strconv.Itoa(i) + "]:" + "data:image/png;base64," + base64Imgs[i]
	}
	return contentMdStr, saveImageBytes
}

func formatBlockQuote(piece parse.Piece, depth int) (string, map[string][]byte) {
	var bqMdString string
	var prefix string = ">"
	for i := 0; i < depth; i++ {
		prefix += ">"
	}
	prefix += " "
	var saveImageBytes map[string][]byte
	bqMdString, saveImageBytes = formatContent(piece.Val.([]parse.Piece), depth+1)
	return prefix + bqMdString + "  \n", saveImageBytes
}

func formatList(li parse.Piece, depth int) (string, map[string][]byte) {
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
	var saveImageBytes map[string][]byte
	listMdString, saveImageBytes = formatContent(li.Val.([]parse.Piece), depth+1)
	return prefix + listMdString + "  \n", saveImageBytes
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

// 图片地址为本身src
func formatImageInline(piece parse.Piece) string {
	return "![" + piece.Attrs["alt"] + "](" + piece.Attrs["src"] + " \"" + piece.Attrs["title"] + "\")"
}

// 图片地址为本地引用
func formatImageFileReferInline(alt string, refName string) string {
	return "![" + alt + "](" + refName + ")"
}

// 图片转成base64并插在原地
func formatImageBase64Inline(piece parse.Piece) string {
	return "![" + piece.Attrs["alt"] + "](data:image/png;base64," + piece.Val.(string) + ")"
}

// 图片地址为markdown内引用（用于base64）
func formatImageRefer(piece parse.Piece, index int) string {
	return "![" + piece.Attrs["alt"] + "][" + strconv.Itoa(index) + "]"
}

func formatLink(piece parse.Piece) string {
	var linkMdStr string = "[" + piece.Val.(string) + "](" + piece.Attrs["href"] + ")"
	return linkMdStr
}
