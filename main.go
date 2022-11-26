package main

import (
	"fmt"
	"os"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
	"github.com/fengxxc/wechatmp2markdown/server"
)

func main() {
	// test.Test1()
	// test.Test2()
	args := os.Args
	if len(args) <= 2 {
		panic("not enough args")
	}
	args1 := args[1]
	args2 := args[2]

	if args1 == "server" {
		// server pattern
		port := args2
		if port == "" {
			port = "8964"
		}
		server.Start(":" + port)
		return
	}

	// cli pattern
	url := args1
	filename := args2
	fmt.Printf("url: %s, filename: %s\n", url, filename)
	var articleStruct parse.Article = parse.ParseFromURL(url)
	format.FormatAndSave(articleStruct, filename)
}
