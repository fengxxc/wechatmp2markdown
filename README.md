# wechatmp2markdown 微信公众号文章转Markdown

## 使用
执行命令：`本程序可执行文件 [url] [filename]`
- `url` 微信公众号文章网页的url
- `filename` makedown文件的保存路径和文件名

例如：windows环境，想把url为`https://mp.weixin.qq.com/s/abcdefg`的文章转成markdown存到 `D:\wechatmp_bak\abc.md` 

则执行： `wechatmp2makrdown_win64.exe https://mp.weixin.qq.com/s/abcdefg D:\wechatmp_bak\abc.md`

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
没装gcc也无妨，去Makefile找对应的命令执行

编译好的文件在`./build`目录下