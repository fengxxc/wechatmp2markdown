package main

import (
	"fmt"
	"os"
	"strings"

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

	// --image=base64 	-ib 保存图片，base64格式，在md文件中（默认为此选项）
	// --image=url 		-iu 只保留图片链接
	// --image=save 	-is 保存图片，最终输出到文件夹
	// --save=zip -sz 		最终打包输出到zip
	imageArgValue := "base64"
	if len(args) > 3 && args[3] != "" {
		if strings.HasPrefix(args[3], "--image=") {
			imageArgValue = args[3][len("--image="):]
		} else if strings.HasPrefix(args[3], "-i") {
			imageArgVal := args[3][len("-i"):]
			switch imageArgVal {
			case "u":
				imageArgValue = "url"
			case "s":
				imageArgValue = "save"
			case "b":
				fallthrough
			default:
				imageArgValue = "base64"
			}
		}
	}

	var imagePolicy parse.ImagePolicy = parse.ImageArgValue2ImagePolicy(imageArgValue)

	// cli pattern
	url := args1
	filename := args2
	fmt.Printf("url: %s, filename: %s\n", url, filename)
	var articleStruct parse.Article = parse.ParseFromURL(url, imagePolicy)
	format.FormatAndSave(articleStruct, filename)
}
