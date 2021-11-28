package parse

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseSection(s *goquery.Selection) Block {
	// fmt.Printf("s.Length() = %d\n", s.Length())
	// fmt.Printf("s.Size() = %d\n", s.Size())
	// var tokens = make([]Token, s.Size())
	var tokens []Token
	s.Children().Each(func(i int, s *goquery.Selection) {
		var t Token
		attr := make(map[string]string)
		if s.Is("span") {
			t = Token{NORMAL_TEXT, s.Text(), nil}
		} else if s.Is("a") {
			attr["href"], _ = s.Attr("href")
			t = Token{LINK, removeBrAndBlank(s.Text()), attr}
		} else if s.Is("img") {
			attr["src"], _ = s.Attr("src")
			t = Token{IMAGE, "", attr}
		} else {
			t = Token{NORMAL_TEXT, s.Text(), nil}
			// TODO
		}
		// fmt.Printf("i = %d\n", i)
		// fmt.Printf("%+v\n", t)
		// tokens[i] = t
		tokens = append(tokens, t)
	})
	return Block{tokens}
}

func ParseFromReader(r io.Reader) Article {
	var article Article
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	var mainContent *goquery.Selection = doc.Find("#img-content")

	// 标题
	title := mainContent.Find("#activity-name").Text()
	fmt.Println(title)
	article.title = title

	// meta 细节待完善
	meta := mainContent.Find("#meta_content").Text()
	meta = removeBrAndBlank(meta)
	fmt.Println(meta)
	article.meta = meta

	// tags 细节待完善
	tags := mainContent.Find("#js_tags").Text()
	tags = removeBrAndBlank(tags)
	fmt.Println(tags)
	article.tags = tags

	// content
	// section[style="line-height: 1.5em;"]>span,a	=> 一般段落（含文本和超链接）
	// p[style="line-height: 1.5em;"]				=> 项目列表（有序/无序）
	// section[style=".*text-align:center"]>img		=> 居中段落（图片）
	content := mainContent.Find("#js_content")
	var sections []Block
	content.Find("section,p").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		// fmt.Println(s.Attr("style"))
		var block Block
		if s.Is("p") {
			block = parseSection(s)
		} else if s.Is("section") {
			block = parseSection(s)
		} else {
			// TODO
		}
		// sections[i] = block
		sections = append(sections, block)
	})
	article.content = sections

	return article
}

func ParseFromHTMLString(s string) Article {
	return ParseFromReader(strings.NewReader(s))
}

func removeBrAndBlank(s string) string {
	regstr := "\\s{2,}"
	reg, _ := regexp.Compile(regstr)
	sb := make([]byte, len(s))
	copy(sb, s)
	spc_index := reg.FindStringIndex(string(sb)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		sb = append(sb[:spc_index[0]+1], sb[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(sb))            //继续在字符串中搜索
	}
	return strings.Replace(string(sb), "\n", " ", -1)
}

var testHTML = `
<div>
    <div id="img-content" class="rich_media_wrp">
        <h1 class="rich_media_title" id="activity-name">终于有人喷我了！</h1>
        <div id="meta_content" class="rich_media_meta_list">
            <span id="copyright_logo" class="wx_tap_link js_wx_tap_highlight rich_media_meta icon_appmsg_tag appmsg_title_tag weui-wa-hotarea" wah-hotarea="click">原创</span>
            <span class="rich_media_meta rich_media_meta_text">
                <span role="link" id="js_author_name" class="wx_tap_link js_wx_tap_highlight weui-wa-hotarea" datarewardsn="" datatimestamp="" datacanreward="0" wah-hotarea="click">闪客sun</span>
            </span>
            <span class="rich_media_meta rich_media_meta_nickname" id="profileBt" wah-hotarea="click">
                <a href="javascript:void(0);" class="wx_tap_link js_wx_tap_highlight weui-wa-hotarea" id="js_name">
                    低并发编程 </a>
                <div id="js_profile_qrcode" aria-hidden="true" class="profile_container" style="display:none;" wah-hotarea="click">
                    <div class="profile_inner">
                        <strong class="profile_nickname">低并发编程</strong>
                        <img class="profile_avatar" id="js_profile_qrcode_img" src="" alt="">
                        <p class="profile_meta">
                            <label class="profile_meta_label">微信号</label>
                            <span class="profile_meta_value">dibingfa</span>
                        </p>
                        <p class="profile_meta">
                            <label class="profile_meta_label">功能介绍</label>
                            <span class="profile_meta_value">战略上藐视技术，战术上重视技术</span>
                        </p>
                    </div>
                    <span class="profile_arrow_wrp" id="js_profile_arrow_wrp">
                        <i class="profile_arrow arrow_out"></i>
                        <i class="profile_arrow arrow_in"></i>
                    </span>
                </div>
            </span>
            <em id="publish_time" class="rich_media_meta rich_media_meta_text">2021-11-26</em>
        </div>
		<div id="js_tags" class="article-tag__list single-tag__wrp js_single js_wx_tap_highlight wx_tap_card" data-len="1" role="link" aria-labelledby="js_article-tag-card__left" aria-describedby="js_article-tag-card__right" wah-hotarea="click">
			<span aria-hidden="true" id="js_article-tag-card__left" class="article-tag-card__left">
				<span class="article-tag-card__title">收录于话题</span>
				<span class="article-tag__item-wrp no-active js_tag" data-url="https://mp.weixin.qq.com/mp/appmsgalbum?__biz=Mzk0MjE3NDE0Ng==&amp;action=getalbum&amp;album_id=1645521656368381958#wechat_redirect" data-tag_id="" data-album_id="1645521656368381958" data-tag_source="4">
					<span class="article-tag__item">#随便聊聊</span>
				</span>
			</span>
			<span aria-hidden="true" id="js_article-tag-card__right" class="article-tag-card__right">49个<span class="weui-hidden_abs">内容</span></span>
		</div>
		<div class="rich_media_content " id="js_content" style="visibility: visible;">
			<p data-mpa-powered-by="yiban.io"><span style="font-size: 16px;letter-spacing: 0.5px;background-color: transparent;caret-color: var(--weui-BRAND);">写公众号一年了，一直盼着能有人喷喷我，今天终于被我碰到了！</span><br></p>
			<section style="line-height: 1.5em;"><span style="letter-spacing: 0.5px;font-size: 16px;">‍‍</span></section>
			<section style="line-height: 1.5em;"><span style="letter-spacing: 0.5px;font-size: 16px;">还是我的一位读者发现的，在推特上，于是分享给了我。<br></span></section>
			<section style="line-height: 1.5em;text-align: center;"><img class="rich_pages wxw-img" data-galleryid="" data-ratio="1.0826306913996628" data-s="300,640" data-src="https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRdNibLb2hf6QnMnWgic4Nm0KhCmicJibxESMoGfbuMrXbQB7lrYFSJPlBeGaJyciaavIBN8NLwESxia7cA/640?wx_fmt=png" data-type="png" data-w="593" style="box-shadow: rgb(210, 210, 210) 0em 0em 0.5em 0px; font-size: 17px; width: 346px !important; height: auto !important; visibility: visible !important;" _width="346px" src="https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRdNibLb2hf6QnMnWgic4Nm0KhCmicJibxESMoGfbuMrXbQB7lrYFSJPlBeGaJyciaavIBN8NLwESxia7cA/640?wx_fmt=png&amp;tp=webp&amp;wxfrom=5&amp;wx_lazy=1&amp;wx_co=1" crossorigin="anonymous" alt="图片" data-fail="0"></section>
			<section style="line-height: 1.5em;">
				<span style="letter-spacing: 0.5px;font-size: 16px;">这是喷我最近的一个新系列，</span>
				<a target="_blank" href="https://mp.weixin.qq.com/mp/appmsgalbum?__biz=Mzk0MjE3NDE0Ng==&amp;action=getalbum&amp;album_id=2123743679373688834#wechat_redirect" textvalue="你管这破玩意叫操作系统源码" linktype="text" imgurl="" imgdata="null" tab="innerlink" data-linktype="2" style="letter-spacing: 0.5px;font-size: 16px;" wah-hotarea="click">
					<span style="letter-spacing: 0.5px;font-size: 16px;">你管这破玩意叫操作系统源码</span>
				</a>
				<span style="letter-spacing: 0.5px;font-size: 16px;">，正愁找不到借口推广一波呢，这不就给我来素材了。<br></span>
			</section>
		</div>
    </div>
</div>
`

func Test() {
	res := ParseFromHTMLString(testHTML)
	fmt.Println("---------------------------------------")
	fmt.Printf("%+v\n", res)
}
