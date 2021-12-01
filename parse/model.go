package parse

type Article struct {
	Title   Piece
	Meta    string
	Tags    string
	Content []Piece
}

type Header struct {
	Level int
	Text  string
}

// go不资瓷泛型可真是难受...
type Value interface{}

type Piece struct {
	Type  PieceType
	Val   Value
	Attrs map[string]string
}

type PieceType int32

const (
	HEADER           PieceType = iota // 标题
	LINK                              // 链接
	NORMAL_TEXT                       // 文字
	BOLD_TEXT                         // 粗体文字
	ITALIC_TEXT                       // 斜体文字
	BOLD_ITALIC_TEXT                  // 粗斜体
	IMAGE                             // 图片
	TABLE                             // 表格
	CODE_INLINE                       // 代码 内联
	CODE_BLOCK                        // 代码 块
	BLOCK_QUOTES                      // 引用
	O_LIST                            // 有序列表
	U_LIST                            // 无序列表
	HR                                // 分隔线
	BR                                // 换行
)
