package test

import (
	"fmt"

	"github.com/fengxxc/wechatmp2markdown/parse"
)

func Test1() {
	res := parse.ParseFromHTMLFile("./test/test1.html")
	fmt.Println("-------------------test1.html-------------------")
	fmt.Printf("%+v\n", res)
}
