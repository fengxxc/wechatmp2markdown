package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
)

func Start(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wechatmpURL := r.FormValue("url")
		if wechatmpURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("param 'url' must not be empty. please put in a wechatmp URL and try again."))
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		var articleStruct parse.Article = parse.ParseFromURL(wechatmpURL)
		title := articleStruct.Title.Val.(string)
		w.Header().Set("Content-Disposition", "attachment; filename="+title+".md")
		var mdString string = format.Format(articleStruct)
		w.Write([]byte(mdString))
		return
	})

	fmt.Printf("wechatmp2markdown server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
