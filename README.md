# wechatmp2markdown 微信公众号文章转Markdown

```
注意：本项目源代码目前只托管在 Github 平台，该平台的logo是一只“直立着、有一只触手的章鱼猫”，若您在其他地方看到本项目，请勿相信。
gggggggggggggggggggggggggggggggggggggggggggggggpppbbbbbbbbbbbbbbbpppggggggggggggggggggggggggggggggggggggggggggggggg
ggggggggggggggggggggggggggggggggggggggpbmhb7                          77hhbbpgggggggggggggggggggggggggggggggggggggg
ggggggggggggggggggggggggggggggggpbh7                                          77hbpgggggggggggggggggggggggggggggggg
gggggggggggggggggggggggggggpph7                                                     7hpgggggggggggggggggggggggggggg
gggggggggggggggggggggggpph                                                              7hpgggggggggggggggggggggggg
ggggggggggggggggggggpmb                                                                     hbggggggggggggggggggggg
gggggggggggggggggppb                                                                           hpgggggggggggggggggg
gggggggggggggggpb                                                                                 hgggggggggggggggg
gggggggggggggpb                                                                                     hpggggggggggggg
gggggggggggpm                                                                                         7pggggggggggg
gggggggggpb               ggggggggggg                                         gggggggggg                hgggggggggg
ggggggggm                 ggggggggggggggg       ggggggggggggggggggg       ggggggggggggggg                 bgggggggg
ggggggpm                 gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                  7ggggggg
gggggm                   gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                    gggggg
ggggm                     ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     qgggg
gggm                      ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggp                      qggg
ggm                      ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                      ggg
gp                     gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     7gg
g                     ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                    jg
p                    gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     g
                     ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                    g
m                    ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                    7
                     ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     
                     gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggm                    
                     ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     
                     iggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg                     g
                      ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggp                     g
g                     7gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggm                     g
g                      7ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggpm                     gg
gg                       bgggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggm                      ggg
ggg                        bpgggggggggggggggggggggggggggggggggggggggggggggggggggggggggpm                        ggg
gggg                         7bpgggggggggggggggggggggggggggggggggggggggggggggggggggpmb                         gggg
ggggg          ggggggg           7hbppgggggggggggggggggggggggggggggggggggggggppbh7                            ggggg
gggggg         7bgggggggg              77hhhbbpgggggggggggggggggggggpbbhhb7                                  gggggg
ggggggg           7bggggggg                  ggggggggggggggggggggggggg                                     gggggggg
ggggggggg            bggggggg               ggggggggggggggggggggggggggg                                   ggggggggg
ggggggggggg           7ggggggggg          gggggggggggggggggggggggggggggg                                ggggggggggg
gggggggggggg           7gggggggggggggggggggggggggggggggggggggggggggggggg                              ggggggggggggg
ggggggggggggggg          bgggggggggggggggggggggggggggggggggggggggggggggg                            ggggggggggggggg
gggggggggggggggggggggggggggggggggggg      gggggggggggggggggggggggggggggg      ggggggggggggggggggggggggggggggggggggg
  _______  __  .___________. __    __   __    __  .______   
 /  _____||  | |           ||  |  |  | |  |  |  | |   _  \  
|  |  __  |  | `---|  |----`|  |__|  | |  |  |  | |  |_)  | 
|  | |_ | |  |     |  |     |   __   | |  |  |  | |   _  <  
|  |__| | |  |     |  |     |  |  |  | |  `--'  | |  |_)  | 
 \______| |__|     |__|     |__|  |__|  \______/  |______/  

```
## 使用
### CLI 模式
通过命令行使用
执行命令：`本程序可执行文件 [url] [filepath] [--image]`
- `url`      微信公众号文章网页的url
- `filepath` makedown文件的保存位置，若该值为目录，则以文章标题作为文件名保存在该目录下；若以`.md`结尾，则以输入的文件名作为文件名保存；`./`为保存到当前目录
- `--image` 可选参数，文章内图片的保存方式，格式为`--image=xxx`，`xxx`为参数值，有三个可供选择（默认值为base64）：
    - `url` 图片引用原src值，它通常在网络上（不推荐，微信哪天把它ban掉就寄了）；
    - `save` 图片存在本地，在与markdown同一个目录中，若为web server模式，则一并打包成zip下载；
    - `base64` 图片编码成base64字符串放在markdown文件内

例如：windows环境，想把url为`https://mp.weixin.qq.com/s/a=1&b=2`的文章（假设文章标题为"gitcode操你妈"）转成markdown存到 `D:\wechatmp_bak`下，文章内的**图片**保存到**本地**

则cmd执行： 
```
wechatmp2makrdown_win64.exe https://mp.weixin.qq.com/s/a=1&b=2 D:\wechatmp_bak --image=save
```

markdown和图片文件将保存在 `D:\wechatmp_bak\gitcode操你妈\` 下

> 在Windows环境下，文件或路径名不能包含以下任何字符：“（双引号）、*（星号）、<（小于）、>（大于）、？（问号）、\（反斜杠）、/（正斜杠）、|（竖线）、：（冒号）。当标题包含以上字符时，本程序将用相似的Unicode字符进行替换，具体替换规则为：  
> ```
> "<" -> "≺"
> ">" -> "≻"
> ":" -> "∶"
> """ -> "“"
> "/" -> "∕"
> "\" -> "∖"
> "|" -> "∣"
> "?" -> "?"
> "*" -> "⁎"
> ```

> 在linux环境下，使用CLI模式，需要为程序赋予可执行与写权限，例如：` chmod +xw wechatmp2markdown-v1.1.8_linux_amd64`

### web server 模式
通过web服务使用

执行命令：`本程序可执行文件 server [port]`
- `port` 监听的端口

当看到 `wechatmp2markdown server listening on :[port]` 时，
打开浏览器（或curl工具）访问：`localhost:[port]?url=[url]&image=[image]`
- `url`   微信公众号文章网页的url
- `image` 可选参数，文章内图片的保存方式，参数值与上文CLI模式的相同

返回的数据即为该文章的markdown文件（若image=save，则返回的是zip格式的压缩包）

例如：windows环境，服务启动并监听8964端口，想把url为`https://mp.weixin.qq.com/s/a=1&b=2`的文章转成markdown并下载，文章内的**图片**保存到**本地**

则cmd执行： `wechatmp2makrdown_win64.exe server 8964`

浏览器访问：`localhost:8964?url=https://mp.weixin.qq.com/s/a=1&b=2&image=save`，
将返回一个zip文件

> 在linux环境下，使用web server模式，需要为程序赋予可执行权限，例如：` chmod +x wechatmp2markdown-v1.1.8_linux_amd64`

## 开发
go 1.20
### 编译
#### Linux or Mac 环境
```
# 编译目标平台: linux
make build-linux
# 编译目标平台: mac
make build-osx
# 编译目标平台: windows 64位
make build-win64
# 编译目标平台: windows 32位
make build-win32
```
#### windows 环境
需gcc，推荐安装tdm64-gcc，`mingw64/bin`加入系统环境变量
```
# 编译目标平台: linux
mingw32-make win-build-linux

# 编译目标平台: mac
mingw32-make win-build-osx

# 编译目标平台: windows 64位
mingw32-make win-build-win64

# 编译目标平台: windows 32位
mingw32-make win-build-win32
```
（没装gcc也无妨，去Makefile找对应的命令执行）

编译好的文件在`./build`目录下

## TODO
- [x] 支持解析表格元素(table tag)

## 最后
抓紧时间窗口。努力记录。黑暗中记得光的模样。
