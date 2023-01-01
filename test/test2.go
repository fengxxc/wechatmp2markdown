package test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fengxxc/wechatmp2markdown/format"
	"github.com/fengxxc/wechatmp2markdown/parse"
)

func Test2() {
	// var articleStruct parse.Article = parse.ParseFromURL("https://mp.weixin.qq.com/s?__biz=MzIzOTU0NTQ0MA==&mid=2247506315&idx=1&sn=1546be4ecece176f669da4eed7076ee2&chksm=e92ae484de5d6d92d93cd68b927fa91e2935a75c9aafc02f294237653ca8a342e8982cabbc1d&cur_album_id=1391790902901014528&scene=189#wechat_redirect")
	var articleStruct parse.Article = parse.ParseFromURL("https://mp.weixin.qq.com/s?__biz=MzU0OTE4MzYzMw==&mid=2247525863&idx=2&sn=d759f98b62f61f3a8312da4ee426c287&chksm=fbb1ec19ccc6650f40c0ef67b47163040c33f9dfe3d6f05bf28d4d823b6f847c09fea046b2eb&scene=132#wechat_redirect", parse.IMAGE_POLICY_BASE64)

	byteArry, _ := json.MarshalIndent(articleStruct, "", "  ")
	// fmt.Println(string(byteArry))
	ioutil.WriteFile("./test/test2_target.json", byteArry, 0644)

	mdString, _ := format.Format((articleStruct))
	ioutil.WriteFile("./test/test2_target.md", []byte(mdString), 0644)
}
