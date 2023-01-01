package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
	"github.com/fengxxc/wechatmp2markdown/util"
)

func Start(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rawQuery := r.URL.RawQuery
		paramsMap := parseParams(rawQuery)

		// url param
		wechatmpURL := paramsMap["url"]
		fmt.Printf("accept url: %s\n", wechatmpURL)
		imageArgValue := paramsMap["image"]
		fmt.Printf("     image: %s\n", imageArgValue)
		imagePolicy := parse.ImageArgValue2ImagePolicy(imageArgValue)

		if wechatmpURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(defHTML))
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		var articleStruct parse.Article = parse.ParseFromURL(wechatmpURL, imagePolicy)
		title := articleStruct.Title.Val.(string)
		mdString, saveImageBytes := format.Format(articleStruct)
		if len(saveImageBytes) > 0 {
			w.Header().Set("Content-Disposition", "attachment; filename="+title+".zip")
			saveImageBytes[title] = []byte(mdString)
			util.HttpDownloadZip(w, saveImageBytes)
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename="+title+".md")
			w.Write([]byte(mdString))
		}
	})

	fmt.Printf("wechatmp2markdown server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

var defHTML string = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>wechatmp2markdown</title>
</head>
<body>
	<h1 style="text-align: center; width: 100%;">wechatmp2markdown</h1>
	<ul style="margin: 0 auto; width: 89%;">
		<li>
			<strong>param 'url' is required.</strong> please put in a wechatmp URL and try again.
		</li>
		<li>
			<strong>param 'image' is optional</strong>, value include: 'url' / 'save' / 'base64'(default)
		</li>
		<li>
			<strong>example:</strong> http://localhost:8964/?url=https://mp.weixin.qq.com/s?__biz=aaaa==&mid=1111&idx=2&sn=bbbb&chksm=cccc&scene=123&image=save
		</li>
	</ul>
</body>
</html>
`

func parseParams(rawQuery string) map[string]string {
	result := make(map[string]string)
	reg := regexp.MustCompile(`(&?image=)([a-z]+)`)
	matcheImage := reg.FindStringSubmatch(rawQuery)
	var urlParamFull string = rawQuery
	if len(matcheImage) > 1 {
		// 有image参数
		imageParamFull := matcheImage[0]
		urlParamFull = strings.Replace(rawQuery, imageParamFull, "", 1)

		if len(matcheImage) > 2 {
			imageParamVal := matcheImage[2]
			result["image"] = imageParamVal
		}
	}
	regUrl := regexp.MustCompile(`(&?url=)(.+)`)
	matcheUrl := regUrl.FindStringSubmatch(urlParamFull)
	if len(matcheUrl) > 2 {
		urlParamVal := matcheUrl[2]
		result["url"] = urlParamVal
	}
	return result
}
