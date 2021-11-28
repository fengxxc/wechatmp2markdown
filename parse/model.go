package parse

type Article struct {
	title   string
	meta    string
	tags    string
	content []Block
}

type Block struct {
	tokens []Token
}

type Token struct {
	ttype TokenType
	text  string
	attrs map[string]string
}

type TokenType int32

const (
	TITLE       TokenType = iota // 标题
	LINK                         // 链接
	NORMAL_TEXT                  // 文字
	STRONG_TEXT                  // 强调文字
	ITALIC_TEXT                  // 斜体文字
	IMAGE                        // 图片
	TABLE                        // 表格
	CODE_INLINE                  // 代码 内联
	CODE_BLOCK                   // 代码 块
	CITE                         // 引用
)
