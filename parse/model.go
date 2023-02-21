package parse

type Article struct {
	Title   Piece
	Meta    []string
	Tags    string
	Content []Piece
}

type Header struct {
	Level int
	Text  string
}

type Value any

type Piece struct {
	Type  PieceType
	Val   Value
	Attrs map[string]string
}

type PieceType int32

const (
	HEADER           PieceType = iota // 0  标题
	LINK                              // 1  链接
	NORMAL_TEXT                       // 2  文字
	BOLD_TEXT                         // 3  粗体文字
	ITALIC_TEXT                       // 4  斜体文字
	BOLD_ITALIC_TEXT                  // 5  粗斜体
	IMAGE                             // 6  图片
	IMAGE_BASE64                      // 7  图片 base64
	TABLE                             // 8  表格
	CODE_INLINE                       // 9  代码 内联
	CODE_BLOCK                        // 10  代码 块
	BLOCK_QUOTES                      // 11 引用
	O_LIST                            // 12 有序列表
	U_LIST                            // 13 无序列表
	HR                                // 14 分隔线
	BR                                // 15 换行
)
