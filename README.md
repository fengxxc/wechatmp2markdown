# wechatmp2markdown 微信公众号文章转Markdown

## 使用
### CLI 模式
通过命令行使用
执行命令：`本程序可执行文件 [url] [filename]`
- `url` 微信公众号文章网页的url
- `filename` makedown文件的保存路径和文件名

例如：windows环境，想把url为`https://mp.weixin.qq.com/s/abcdefg`的文章转成markdown存到 `D:\wechatmp_bak\abc.md` 

则cmd执行： `wechatmp2makrdown_win64.exe https://mp.weixin.qq.com/s/abcdefg D:\wechatmp_bak\abc.md`
### web server 模式
通过web服务使用
执行命令：`本程序可执行文件 server [port]`
- `port` 监听的端口

当看到 `wechatmp2markdown server listening on :[port]` 时，
打开浏览器（或curl工具）访问：`localhost:[port]?url=[url]`
- `url` 微信公众号文章网页的url

返回的数据即为该文章的markdown文件

例如：windows环境，服务启动并监听8080端口，想把url为`https://mp.weixin.qq.com/s/abcdefg`的文章转成markdown并下载

则cmd执行： `wechatmp2makrdown_win64.exe server 8080`

浏览器访问：`localhost:8080?url=https://mp.weixin.qq.com/s/abcdefg`，

## 开发
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

## 最后
抓紧时间窗口。努力记录。黑暗中记得光的模样。