package test

import (
	"fmt"
	"io/ioutil"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
)

func Test1() {
	var articleStruct parse.Article = parse.ParseFromHTMLFile("./test/test1.html")
	fmt.Println("-------------------test1.html parse-------------------")
	fmt.Printf("%+v\n", articleStruct)

	fmt.Println("-------------------test1.html format-------------------")
	var mdString string = format.Format(articleStruct)
	fmt.Print(mdString)
	ioutil.WriteFile("./test/test1_target.md", []byte(mdString), 0644)
}
