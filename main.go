package main

import (
	"fmt"
	"os"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
)

func main() {
	// test.Test1()
	// test.Test2()
	args := os.Args
	if len(args) <= 2 {
		panic("not enough args")
	}
	url := args[1]
	filename := args[2]
	fmt.Printf("url: %s, filename: %s\n", url, filename)
	var articleStruct parse.Article = parse.ParseFromURL(url)
	format.FormatAndSave(articleStruct, filename)
}
