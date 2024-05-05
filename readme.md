## ZxwyWebSite/cr-go-sdk
### 简介
+ Cloudreve Golang SDK 接口封装
<!-- + "面向对象，多站点，多用户，兼容多版本接口" -->
+ 测试阶段，存在不稳定因素，仅供参考，不建议用于生产环境
+ 基于社区版 V3.8.3 开发，这是V3系列的最后一个版本，请尽量更新到此版本以保证稳定运行
+ 目前只是基于Api编写了接口调用方法，具体功能还需另外实现
+ 原生方法以Api路径命名，相同路径根据操作命名，封装方法以Sdk开头
+ example目录里有一个演示程序

### 使用
0. 安装依赖
   ```sh
   go get -u github.com/ZxwyWebSite/cr-go-sdk
   ```
1. 创建站点
   ```go
   site, _ := cr.NewSite(`https://cloudreveplus-demo.onrender.com/`, cr.ApiV353)
   ```
2. 登录账号
   ```go
   // 添加账户信息
   site.Users = &cr.UserObj{
   	Mail: `admin@cloudreve.org`,
   	Pass: `CloudrevePlusDemo`,
   	Cookie: cr.ParseCookie(`cloudreve-session=xxx; Path=/; Expires=Sat, 11 May 2024 09:05:02 GMT; Max-Age=604800; HttpOnly`),
   }
   // 初始化站点数据
   err = site.SdkInit()
   // 登录账号
   err = site.SdkLogin()
   ```
3. 执行操作
   ```go
   // 列出文件
   list, err := site.Directory(`/`)
   if err != nil {
   	panic(err)
   }
   fmt.Printf("%#v\n", list)
   // 下载文件
   link, err := site.FileDownload(`XBUl`)
   if err != nil {
   	panic(err)
   }
   fmt.Println(*link)
   // 上传文件
   file, _ := os.Open(`.outdated/test.mp4`)
   err = site.SdkUpload(`/`, file, `test_0.mp4`)
   ```
<!--3. 操作文件
   ```go
   dir, _ := user.File.List(`/`)
   dir.Range(func (f *cr.FileObj) bool {
       if !f.IsDir {
           link, _ := f.Download()
           fmt.Println(f.Name, link)
       }
       return true
   })
   ```-->

### 结构
+ models|pkg|serializer|service/ 参数定义
+ uploader/ 上传组件
+ api.go 原版接口定义
+ sdk.go 扩展功能封装
+ site.go 站点对象
+ user.go 账户对象
+ util.go 实用工具

### 其它
+ EOF
<!-- + CloudreamProject云梦企划 -->