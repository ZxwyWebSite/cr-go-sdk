package cr

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unsafe"

	json "github.com/ZxwyProject/zson"
	"github.com/ZxwyWebSite/cr-go-sdk/pkg/payment"
	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/aria2"
	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/setting"
	"github.com/ZxwyWebSite/cr-go-sdk/service/share"
	"github.com/ZxwyWebSite/cr-go-sdk/service/user"
	"github.com/ZxwyWebSite/cr-go-sdk/service/vas"
)

// 获取接口相对路径 [回调] [结果]
func (c *SiteObj) api(f func(*strings.Builder)) string {
	var b strings.Builder
	/*if n != 0 {
		b.Grow(len(c.Addr) + 7 + n)
	}*/
	b.WriteString(c.Addr)
	b.WriteString(`api/v3/`)
	f(&b)
	return b.String()
}

// 创建网络请求 [方法,地址,数据,登录] [请求,错误]
func (c *SiteObj) newRequest(method string, url string, body io.Reader, sess bool) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err == nil {
		req.Header[`User-Agent`] = []string{Cr_UserAgent}
		req.Header[`Accept`] = []string{Cr_Accept}
		if sess {
			if c.Users != nil {
				if c.Users.Cookie != nil {
					if c.Users.Cookie.RawExpires != `` {
						if time.Now().Before(c.Users.Cookie.Expires) {
							if c.Users.Cookie.Raw != `` {
								req.Header[`Cookie`] = []string{c.Users.Cookie.Raw}
								// Cr_Format(`[debug] rawcookie:`, c.Users.Cookie.Raw)
							} else {
								req.Header[`Cookie`] = []string{c.Users.Cookie.String()}
								// Cr_Format(`[debug] addcookie:`, req.Header[`Cookie`][0])
							}
						}
					}
				}
			}
		}
	}
	return req, err
}

// 发送网络请求 [方法,路径,参数,映射] [错误]
func (c *SiteObj) fetch(method string, uri func(b *strings.Builder), s any, out any) error {
	var body io.Reader
	if s != nil {
		// 构建请求体
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(s); err != nil {
			return err
		}
		body = &buf
	}
	req, err := c.newRequest(method, c.api(uri), body, true)
	if err != nil {
		return err
	}
	if s != nil {
		req.Header[`Content-Type`] = []string{`application/json`}
	}
	res, err := Cr_Client.Do(req)
	if err != nil {
		return err
	}
	if c.Users != nil {
		/*if sess := res.Cookies(); len(sess) >= 1 {
			c.Users.Cookie = sess[0]
			if Cr_Debug {
				Cr_Format(`[debug] sessupd:`, sess[0].String())
			}
		}*/
		for _, s := range res.Cookies() {
			if s.Name == `cloudreve-session` {
				c.Users.Cookie = s
				if Cr_Debug {
					Cr_Format(`[debug] sessupd:`, s.String())
				}
			}
		}
	}
	// 注：由于Golang泛型必须先初始化，而初始化后的类型又无法直接断言，故无法在此处检查返回状态码
	// err: cannot use generic type serializer.Response[T any] without instantiation
	// 已暂时通过添加检测接口解决
	if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
		if o, ok := out.(interface{ Err() error }); ok {
			err = o.Err()
		}
	}
	res.Body.Close()
	return err
}

// Site 全局设置相关

// 获取版本号 [] [版本号(3.8.3+1.1-plus),错误]
func (c *SiteObj) SitePing() (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`site/ping`) }, nil, &out)
	return out.Data, err
}

// 获取验证码 [] [base64编码后的图片(data:image/png;base64,),错误]
func (c *SiteObj) SiteCaptcha() (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`site/captcha`) }, nil, &out)
	return out.Data, err
}

// 站点全局配置 [] [配置信息,错误]
func (c *SiteObj) SiteConfig() (*serializer.SiteConfig, error) {
	var out serializer.Response[*serializer.SiteConfig]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`site/config`) }, nil, &out)
	return out.Data, err
}

// 获取 VOL 密钥 (Pro) [] [VOL密钥,错误]
func (c *SiteObj) SiteVol() (*serializer.VolResponse, error) {
	var out serializer.Response[*serializer.VolResponse]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`site/vol`) }, nil, &out)
	return out.Data, err
}

// User 用户相关路由

// 用户登录 [登录凭证] [用户详情,错误]
func (c *SiteObj) UserSession(s *user.LoginInfo) (*serializer.User, error) {
	// if c.Config.LoginCaptcha {
	// 	return nil, errors.New(`暂不支持开启验证码的站点`)
	// }
	var out serializer.Response[*serializer.User]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`user/session`) }, s, &out)
	return out.Data, err
}

// 用户注册 [用户信息] [错误]
func (c *SiteObj) UserReg(s *user.LoginInfo) error {
	// var buf bytes.Buffer
	// enc := json.NewEncoder(&buf)
	// enc.SetEscapeHTML(false)
	// err := enc.Encode(s)
	// if err != nil {
	// 	return err
	// }
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`user`) }, s, &out)
	return err
}

// 用户2FA验证 [] [用户详情,错误]
func (c *SiteObj) User2FA(s *user.Enable2FA) (*serializer.User, error) {
	var out serializer.Response[*serializer.User]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`user/2fa`) }, s, &out)
	return out.Data, err
}

// 发送密码重设邮件 [] []
func (c *SiteObj) UserResetSend(s *user.UserResetEmailService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`user/reset`) }, s, &out)
	return err
}

// 通过邮件里的链接重设密码 [] []
func (c *SiteObj) UserResetSubmit(s *user.UserResetService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPatch, func(b *strings.Builder) { b.WriteString(`user/reset`) }, s, &out)
	return err
}

// 邮件激活 [验证码] [用户邮箱,错误]
func (c *SiteObj) UserActivate(id string) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`user/activate/`)
		b.WriteString(id)
	}, nil, &out)
	return out.Data, err
}

// 初始化QQ登录 [] [登录页地址,错误]
func (c *SiteObj) UserQQ() (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`user/qq`) }, nil, &out)
	return out.Data, err
}

// WebAuthn登陆初始化 [用户名] [(外部资源),错误]
func (c *SiteObj) UserAuthn(username string) (any, error) {
	var out serializer.Response[any]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`user/authn/`)
		b.WriteString(username)
	}, nil, &out)
	return out.Data, err
}

// WebAuthn登陆 [用户名] [用户详情,错误]
func (c *SiteObj) UserAuthnFinish(username string) (*serializer.User, error) {
	var out serializer.Response[*serializer.User]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`user/authn/finish/`)
		b.WriteString(username)
	}, nil, &out)
	return out.Data, err
}

// 获取用户主页展示用分享 [用户id,查询参数(page=0&type=default)] [分享列表,错误]
func (c *SiteObj) UserProfile(id string, q *share.ShareUserGetService) (*serializer.ShareList, error) {
	var b strings.Builder
	b.WriteString(`user/profile/`)
	b.WriteString(id)
	if q.Page != 0 {
		b.WriteString(`?page=`)
		b.WriteString(strconv.FormatUint(uint64(q.Page), 10))
	} else {
		return nil, errors.New(`q.Page cannot be empty`)
	}
	if q.Type != `` {
		b.WriteString(`&type=`)
		b.WriteString(q.Type)
	} else {
		return nil, errors.New(`q.Type cannot be empty`)
	}
	var out serializer.Response[*serializer.ShareList]
	err := c.fetch(http.MethodGet, func(s *strings.Builder) { s.WriteString(b.String()) }, nil, &out)
	return out.Data, err
}

// 获取用户头像 [用户id,图片大小(l|s)] [重定向链接]
func (c *SiteObj) UserAvatar(id, size string) string {
	// var out serializer.Response[*string]
	// err := c.fetch(http.MethodGet, b.String(), nil, &out)
	// return out.Data, err
	return c.api(func(b *strings.Builder) {
		b.WriteString(`user/avatar/`)
		b.WriteString(id)
		b.WriteByte('/')
		b.WriteString(size)
	})
}

// ↑ 需要携带签名验证的 ↓

// 文件外链（直接输出文件数据） [文件id,文件名称] [下载地址]
/*func (c *SiteObj) FileGet(id, name string) (string, error) {
	var b strings.Builder
	b.WriteString(`file/get/`)
	b.WriteString(id)
	b.WriteByte('/')
	b.WriteString(name)
	return b.String(), nil
}*/

// 文件外链(301跳转) [文件id,文件名称] [下载地址]
/*func (c *SiteObj) FileSource(id, name string) (string, error) {
	var b strings.Builder
	b.WriteString(`file/source/`)
	b.WriteString(id)
	b.WriteByte('/')
	b.WriteString(name)
	return b.String(), nil
}*/

// 下载文件 [文件id,文件名称] [下载地址]
/*func (c *SiteObj) FileDownload(id string) (string, error) {
	return `file/download/` + id, nil
}*/

// 打包并下载文件 [会话id(sessionID)] [下载地址]
/*func (c *SiteObj) FileArchive(sid string) (string, error) {
	var b strings.Builder
	b.WriteString(`file/archive/`)
	b.WriteString(sid)
	b.WriteString(`/archive.zip`)
	return b.String(), nil
}*/

// 复制用户会话 [用户id] [错误]
/*func (c *SiteObj) UserSessionCopy(id string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodGet, `user/session/copy/`+id, nil, &out)
	return err
}*/

// ↑ 从机的 RPC 通信 ↓

// ...TODO

// Callback 回调接口

// OneDrive文件上传完成 [会话id] [错误]
func (c *SiteObj) CallbackOneDriveFinish(sid string) error {
	// inline
	req, err := c.newRequest(http.MethodPost, c.api(func(b *strings.Builder) {
		b.WriteString(`callback/onedrive/finish/`)
		b.WriteString(sid)
	}), strings.NewReader(`{}`), true)
	if err != nil {
		return err
	}
	req.Header[`Content-Type`] = []string{`application/x-www-form-urlencoded`}
	res, err := Cr_Client.Do(req)
	if err != nil {
		return err
	}
	var out serializer.Response[struct{}]
	if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
		err = out.Err()
	}
	res.Body.Close()
	return err
}

// AWS S3策略上传回调 [会话id] [错误]
func (c *SiteObj) CallbackS3(sid string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`callback/s3/`)
		b.WriteString(sid)
	}, nil, &out)
	return err
}

// ...TODO

// Share 分享相关

// 获取分享 [分享id,分享码?] [分享信息,错误]
func (c *SiteObj) ShareInfo(id, password string) (*serializer.Share, error) {
	var out serializer.Response[*serializer.Share]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`share/info/`)
		b.WriteString(id)
		if password != `` {
			b.WriteString(`?password=`)
			b.WriteString(url.QueryEscape(password))
		}
	}, nil, &out)
	return out.Data, err
}

// 创建文件下载会话 [分享id,文件路径] [下载链接,错误]
func (c *SiteObj) ShareDownload(id, path string) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) {
		b.WriteString(`share/download/`)
		b.WriteString(id)
		if path != `` {
			b.WriteString(`?path=`)
			b.WriteString(url.QueryEscape(path))
		}
	}, nil, &out)
	if err == nil {
		// 兼容 v3.8.0 破坏性变更 本机存储路径拼接
		if (*out.Data)[0] == '/' {
			*out.Data = c.Addr + (*out.Data)[1:]
		}
	}
	return out.Data, err
}

// 预览分享文件 [分享id,文件路径] [重定向链接(需携带Cookie访问)]
func (c *SiteObj) SharePreview(id, path string) string {
	return c.api(func(b *strings.Builder) {
		b.WriteString(`share/preview/`)
		b.WriteString(id)
		if path != `` {
			b.WriteString(`?path=`)
			b.WriteString(url.QueryEscape(path))
		}
	})
}

// 取得Office文档预览地址 [分享id,文件路径] [预览会话,错误]
func (c *SiteObj) ShareDoc(id, path string) (*serializer.DocPreviewSession, error) {
	var out serializer.Response[*serializer.DocPreviewSession]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`share/doc/`)
		b.WriteString(id)
		if path != `` {
			b.WriteString(`?path=`)
			b.WriteString(url.QueryEscape(path))
		}
	}, nil, &out)
	return out.Data, err
}

// 获取文本文件内容 [分享id,文件路径] [重定向链接,错误]
func (c *SiteObj) ShareContent(id, path string) (string, error) {
	// inline
	req, err := c.newRequest(http.MethodGet, c.api(func(b *strings.Builder) {
		b.WriteString(`share/content/`)
		b.WriteString(id)
		if path != `` {
			b.WriteString(`?path=`)
			b.WriteString(url.QueryEscape(path))
		}
	}), nil, true)
	if err != nil {
		return ``, err
	}
	res, err := Cr_Client.Do(req)
	if err != nil {
		return ``, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		return ``, err
	}
	res.Body.Close()
	return unsafe.String(unsafe.SliceData(data), len(data)), nil
}

// 分享目录列文件 [分享id,文件路径] [目录列表,错误]
func (c *SiteObj) ShareList(id, path string) (*serializer.ObjectList, error) {
	var out serializer.Response[*serializer.ObjectList]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`share/list/`)
		b.WriteString(id)
		if path != `` {
			b.WriteByte('/')
			b.WriteString(url.PathEscape(path))
		}
	}, nil, &out)
	return out.Data, err
}

// 分享目录搜索 [分享id,搜索方式,关键词] [搜索结果,错误]
func (c *SiteObj) ShareSearch(id, Type, keywords string) (*explorer.SearchResult, error) {
	var out serializer.Response[*explorer.SearchResult]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`share/search/`)
		b.WriteString(id)
		b.WriteByte('/')
		b.WriteString(Type)
		b.WriteByte('/')
		b.WriteString(keywords)
	}, nil, &out)
	return out.Data, err
}

// 归档打包下载 [分享id,打包参数] [下载链接,错误]
func (c *SiteObj) ShareArchive(id string, s *share.ArchiveService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`share/archive/`)
		b.WriteString(id)
	}, s, &out)
	return out.Data, err
}

// 获取README文本文件内容 [分享id,文件路径] [下载链接,错误]
func (c *SiteObj) ShareReadme(id, path string) (string, error) {
	var b strings.Builder
	b.WriteString(c.Addr)
	b.WriteString(`api/v3/`)
	b.WriteString(`share/readme/`)
	b.WriteString(id)
	if path != `` {
		b.WriteString(`?path=`)
		b.WriteString(path)
	} else {
		return ``, errors.New(`path cannot be empty`)
	}
	// var out serializer.Response[struct{}]
	// err := c.fetch(http.MethodGet, func(s *strings.Builder) { s.WriteString(b.String()) }, nil, &out)
	return b.String(), nil
}

// 获取缩略图 [分享id,文件名称,文件路径] [下载链接,错误]
func (c *SiteObj) ShareThumb(id, file, path string) (string, error) {
	var b strings.Builder
	b.WriteString(c.Addr)
	b.WriteString(`api/v3/`)
	b.WriteString(`share/thumb/`)
	b.WriteString(id)
	b.WriteByte('/')
	b.WriteString(file)
	if path != `` {
		b.WriteString(`?path=`)
		b.WriteString(path)
	} else {
		return ``, errors.New(`path cannot be empty`)
	}
	// var out serializer.Response[struct{}]
	// err := c.fetch(http.MethodGet, func(s *strings.Builder) { s.WriteString(b.String()) }, nil, &out)
	return b.String(), nil
}

// 举报分享 (Pro) [分享,举报参数] [错误]
func (c *SiteObj) ShareReport(id string, s *share.ShareReportService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`share/report/`)
		b.WriteString(id)
	}, s, &out)
	return err
}

// 搜索公共分享 [搜索参数(page=1&order_by=created_at&order=DESC&keywords=test)] [搜索结果,错误]
func (c *SiteObj) ShareSearchPub(q *share.ShareListService) (*serializer.ShareList, error) {
	var b strings.Builder
	b.WriteString(`share/search`)
	if q.Page != 0 {
		b.WriteString(`?page=`)
		b.WriteString(strconv.FormatUint(uint64(q.Page), 10))
	} else {
		return nil, errors.New(`q.Page cannot be empty`)
	}
	if q.OrderBy != `` {
		b.WriteString(`&order_by=`)
		b.WriteString(q.OrderBy)
	} else {
		return nil, errors.New(`q.OrderBy cannot be empty`)
	}
	if q.Order != `` {
		b.WriteString(`&order=`)
		b.WriteString(q.Order)
	} else {
		return nil, errors.New(`q.Order cannot be empty`)
	}
	if q.Keywords != `` {
		b.WriteString(`&keywords=`)
		b.WriteString(q.Keywords)
	} else {
		return nil, errors.New(`q.Keywords cannot be empty`)
	}
	var out serializer.Response[*serializer.ShareList]
	err := c.fetch(http.MethodGet, func(s *strings.Builder) { s.WriteString(b.String()) }, nil, &out)
	return out.Data, err
}

// ↑ Office Wopi ↓

// ...TODO

// 需要登录保护的

// ↑ 管理 ↓

// ...TODO

// User 用户

// 当前登录用户信息 [] [用户详情,错误]
func (c *SiteObj) UserMe() (*serializer.User, error) {
	var out serializer.Response[*serializer.User]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/me`) }, nil, &out)
	return out.Data, err
}

// 存储信息 [] [存储信息,错误]
func (c *SiteObj) UserStorage() (*serializer.Storage, error) {
	var out serializer.Response[*serializer.Storage]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/storage`) }, nil, &out)
	return out.Data, err
}

// 退出登录 [] [错误]
func (c *SiteObj) UserSessionDel() error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) { b.WriteString(`user/session`) }, nil, &out)
	return err
}

// 生成用于复制客户端会话的临时 URL，用于为移动应用程序添加帐户。 [] [链接,错误]
func (c *SiteObj) UserSessionGen() (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/session`) }, nil, &out)
	return out.Data, err
}

// WebAuthn 注册相关

// ...TODO

// ↑ 用户设置 ↓

// 获取用户可选存储策略 (Pro) [] [存储策略列表,错误]
func (c *SiteObj) UserSettingPolicies() (*[]serializer.PolicyOptions, error) {
	var out serializer.Response[*[]serializer.PolicyOptions]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/setting/policies`) }, nil, &out)
	return out.Data, err
}

// 获取用户可选节点 [] [离线下载节点列表,错误]
func (c *SiteObj) UserSettingNodes() (*[]serializer.NodeOptions, error) {
	var out serializer.Response[*[]serializer.NodeOptions]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/setting/nodes`) }, nil, &out)
	return out.Data, err
}

// 任务队列 [分页] [任务列表,错误]
func (c *SiteObj) UserSettingTasks(page int) (*serializer.TaskList, error) {
	var b strings.Builder
	b.WriteString(`user/setting/tasks`)
	if page != 0 {
		b.WriteString(`?page=`)
		b.WriteString(strconv.Itoa(page))
	} else {
		return nil, errors.New(`page cannot be empty`)
	}
	var out serializer.Response[*serializer.TaskList]
	err := c.fetch(http.MethodGet, func(s *strings.Builder) { s.WriteString(b.String()) }, nil, &out)
	return out.Data, err
}

// 获取当前用户设定 [] [用户设定,错误]
func (c *SiteObj) UserSetting() (*user.Settings, error) {
	var out serializer.Response[*user.Settings]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/setting`) }, nil, &out)
	return out.Data, err
}

// 从文件上传头像 [文件,名称(akari.jpg)<扩展名非常重要>] [错误]
func (c *SiteObj) UserSettingAvatarUpd(file []byte, name string) error {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	pw, err := mw.CreateFormFile(`avatar`, name)
	if err != nil {
		return err
	}
	pw.Write(file)
	// h := make(textproto.MIMEHeader)
	// var b strings.Builder
	// b.WriteString(`form-data; name="`)
	// b.WriteString(`avatar`)
	// b.WriteString(`"; filename="`)
	// b.WriteString(name)
	// b.WriteByte('"')
	// h.Set("Content-Disposition", b.String())
	// h.Set("Content-Type", "image/jpeg")
	// mw.CreatePart(h)
	mw.Close()
	var out serializer.Response[struct{}]
	// inline
	err = func() error {
		req, err := http.NewRequest(http.MethodPost, c.Addr+`api/v3/`+`user/setting/avatar`, buf)
		if err != nil {
			return err
		}
		req.Header[`User-Agent`] = []string{Cr_UserAgent}
		req.Header[`Content-Type`] = []string{mw.FormDataContentType()}
		// req.Header.Set(`Cookie`, c.Users.Sess)
		req.Header[`Cookie`] = []string{c.Users.Cookie.String()}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
			err = out.Err()
		}
		res.Body.Close()
		return err
	}()
	return err
}

// 设定为Gravatar头像 [] [错误]
func (c *SiteObj) UserSettingAvatar() error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) { b.WriteString(`user/setting/avatar`) }, nil, &out)
	return err
}

// 更改用户设定 [服务(nick√|vip|qq|policy|homepage|password|2fa|authn|theme),设定] [错误]
func (c *SiteObj) UserSettingUpd(option string, s *user.UpdOption) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPatch, func(b *strings.Builder) {
		b.WriteString(`user/setting/`)
		b.WriteString(option)
	}, s, &out)
	return err
}

// 获得二步验证初始化信息 [] [密钥,错误]
func (c *SiteObj) UserSetting2FA() (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`user/setting/2fa`) }, nil, &out)
	return out.Data, err
}

// File 文件操作

// ↑ 上传 ↓

// 文件上传 (本地) [会话id,分片数,文件,大小,类型] [错误]
func (c *SiteObj) FileUploadPut(sid string, index int, file io.Reader, size uint64, mime string) error {
	// inline
	req, err := c.newRequest(http.MethodPost, c.api(func(b *strings.Builder) {
		b.WriteString(`file/upload/`)
		b.WriteString(sid)
		b.WriteByte('/')
		b.WriteString(strconv.Itoa(index))
	}), file, true)
	if err != nil {
		return err
	}
	req.ContentLength = int64(size) // req.Header[`Content-Length`] = []string{strconv.FormatUint(size, 10)}
	req.Header[`Content-Type`] = []string{mime}
	res, err := Cr_Client.Do(req)
	if err != nil {
		return err
	}
	var out serializer.Response[struct{}]
	if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
		err = out.Err()
	}
	res.Body.Close()
	return err
}

// 创建上传会话 [文件信息] [上传凭证,错误]
func (c *SiteObj) FileUploadNew(s *explorer.CreateUploadSessionService) (*serializer.UploadCredential, error) {
	var out serializer.Response[*serializer.UploadCredential]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) { b.WriteString(`file/upload`) }, s, &out)
	return out.Data, err
}

// 删除给定上传会话 [会话id] [错误]
func (c *SiteObj) FileUploadDel(sid string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`file/upload/`)
		b.WriteString(sid)
	}, nil, &out)
	return err
}

// 删除全部上传会话 [] [错误]
func (c *SiteObj) FileUploadDelAll() error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) { b.WriteString(`file/upload`) }, nil, &out)
	return err
}

// 更新文件 [文件id,文件内容] [错误]
func (c *SiteObj) FileUpdate(id string, file []byte) error {
	var b strings.Builder
	b.WriteString(c.Addr)
	b.WriteString(`api/v3/`)
	b.WriteString(`file/update/`)
	b.WriteString(id)
	var out serializer.Response[struct{}]
	// inline
	req, err := http.NewRequest(http.MethodPut, b.String(), bytes.NewReader(file))
	if err != nil {
		return err
	}
	req.Header[`User-Agent`] = []string{Cr_UserAgent}
	req.Header[`Content-Length`] = []string{strconv.Itoa(len(file))}
	req.Header[`Content-Type`] = []string{`application/x-www-form-urlencoded`}
	// req.Header.Set(`Cookie`, c.Users.Sess)
	req.Header[`Cookie`] = []string{c.Users.Cookie.String()}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
		err = out.Err()
	}
	res.Body.Close()
	return err
}

// 创建空白文件 [路径参数] [错误]
func (c *SiteObj) FileCreate(s *explorer.SingleFileService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/create`) }, s, &out)
	return err
}

// 创建文件下载会话 [文件id] [下载链接,错误]
func (c *SiteObj) FileDownload(id string) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) {
		b.WriteString(`file/download/`)
		b.WriteString(id)
	}, nil, &out)
	if err == nil {
		// 兼容 v3.8.0 破坏性变更 本机存储路径拼接
		if (*out.Data)[0] == '/' {
			*out.Data = c.Addr + (*out.Data)[1:]
		}
	}
	return out.Data, err
}

// 预览文件 (重定向) [文件id] [链接]
func (c *SiteObj) FilePreview(id string) string {
	var b strings.Builder
	b.WriteString(c.Addr)
	b.WriteString(`api/v3/`)
	b.WriteString(`file/preview/`)
	b.WriteString(id)
	return b.String()
}

// 获取文本文件内容 (重定向) [文件id] [链接]
func (c *SiteObj) FileContent(id string) string {
	return c.api(func(b *strings.Builder) {
		b.WriteString(`file/content/`)
		b.WriteString(id)
	})
}

// 取得Office文档预览地址 [文件id] [预览会话,错误]
func (c *SiteObj) FileDoc(id string) (*serializer.DocPreviewSession, error) {
	var out serializer.Response[*serializer.DocPreviewSession]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`file/doc/`)
		b.WriteString(id)
	}, nil, &out)
	return out.Data, err
}

// 获取缩略图 (重定向) [文件id] [链接]
func (c *SiteObj) FileThumb(id string) string {
	return c.api(func(b *strings.Builder) {
		b.WriteString(`file/thumb/`)
		b.WriteString(id)
	})
}

// 取得文件外链 [文件列表] [外链,错误]
func (c *SiteObj) FileSource(s *explorer.ItemIDService) (*[]serializer.Sources, error) {
	var out serializer.Response[*[]serializer.Sources]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/source`) }, s, &out)
	return out.Data, err
}

// 打包要下载的文件 [文件列表] [结果,错误]
func (c *SiteObj) FileArchive(s *explorer.ItemIDService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/archive`) }, s, &out)
	return out.Data, err
}

// 创建文件压缩任务 [任务参数] [错误]
func (c *SiteObj) FileCompress(s *explorer.ItemCompressService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/compress`) }, s, &out)
	return err
}

// 创建文件解压缩任务 [任务参数] [错误]
func (c *SiteObj) FileDecompress(s *explorer.ItemDecompressService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/decompress`) }, s, &out)
	return err
}

// 创建文件转移任务 [任务参数] [错误]
func (c *SiteObj) FileRelocate(s *explorer.ItemRelocateService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`file/relocate`) }, s, &out)
	return err
}

// 搜索文件 [查询参数(keywords|image|video|audio|doc|tag)] [搜索结果,错误]
func (c *SiteObj) FileSearch(q *explorer.ItemSearchService) (*explorer.SearchResult, error) {
	var out serializer.Response[*explorer.SearchResult]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`file/search/`)
		b.WriteString(q.Type)
		b.WriteByte('/')
		if q.Keywords == `` {
			b.WriteString(`internal`)
		} else {
			b.WriteString(q.Keywords)
		}
		if q.Path != `` {
			b.WriteString(`?path=`)
			b.WriteString(url.PathEscape(q.Path))
		}
	}, nil, &out)
	return out.Data, err
}

// Aria2 离线下载任务

// 创建URL下载任务 [任务参数] [结果,错误]
func (c *SiteObj) Aria2Url(s *aria2.BatchAddURLService) (*[]serializer.Response[struct{}], error) {
	var out serializer.Response[*[]serializer.Response[struct{}]]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`aria2/url`) }, s, &out)
	return out.Data, err
}

// 创建种子下载任务 [文件id,任务参数] [错误]
func (c *SiteObj) Aria2Torrent(id string, s *aria2.AddURLService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`aria2/torrent/`)
		b.WriteString(id)
	}, s, &out)
	return err
}

// 重新选择要下载的文件 [任务id,选择参数] [错误]
func (c *SiteObj) Aria2Select(gid string, s *aria2.SelectFileService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) {
		b.WriteString(`aria2/select/`)
		b.WriteString(gid)
	}, s, &out)
	return err
}

// 取消或删除下载任务 [任务id] [错误]
func (c *SiteObj) Aria2Task(gid string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`aria2/task/`)
		b.WriteString(gid)
	}, nil, &out)
	return err
}

// 获取正在下载中的任务 [查询参数] [列表条目,错误]
func (c *SiteObj) Aria2Downloading(q *aria2.DownloadListService) (*[]serializer.DownloadListResponse, error) {
	var out serializer.Response[*[]serializer.DownloadListResponse]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`aria2/downloading`)
		b.WriteString(`?page=`)
		b.WriteString(strconv.FormatUint(uint64(q.Page), 10))
	}, nil, &out)
	return out.Data, err
}

// 获取已完成的任务 [查询参数] [列表条目,错误]
func (c *SiteObj) Aria2Finished(q *aria2.DownloadListService) (*[]serializer.FinishedListResponse, error) {
	var out serializer.Response[*[]serializer.FinishedListResponse]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`aria2/finished`)
		b.WriteString(`?page=`)
		b.WriteString(strconv.FormatUint(uint64(q.Page), 10))
	}, nil, &out)
	return out.Data, err
}

// Directory 目录

// 创建目录 [路径参数] [错误]
func (c *SiteObj) DirectoryNew(s *explorer.DirectoryService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPut, func(b *strings.Builder) { b.WriteString(`directory`) }, s, &out)
	return err
}

// 列出目录下内容 [路径] [错误]
func (c *SiteObj) Directory(path string) (*serializer.ObjectList, error) {
	var out serializer.Response[*serializer.ObjectList]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`directory`)
		b.WriteString(path)
	}, nil, &out)
	return out.Data, err
}

// Object 对象，文件和目录的抽象

// 删除对象 [文件列表] [错误]
func (c *SiteObj) ObjectDel(s *explorer.ItemIDService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) { b.WriteString(`object`) }, s, &out)
	return err
}

// 移动对象 [文件列表] [错误]
func (c *SiteObj) ObjectMov(s *explorer.ItemMoveService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPatch, func(b *strings.Builder) { b.WriteString(`object`) }, s, &out)
	return err
}

// 复制对象 [文件列表] [错误]
func (c *SiteObj) ObjectCopy(s *explorer.ItemMoveService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`object/copy`) }, s, &out)
	return err
}

// 重命名对象 [文件列表] [错误]
func (c *SiteObj) ObjectRename(s *explorer.ItemRenameService) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`object/rename`) }, s, &out)
	return err
}

// 获取对象属性 [查询参数] [属性信息,错误]
func (c *SiteObj) ObjectProperty(q *explorer.ItemPropertyService) (*serializer.ObjectProps, error) {
	var out serializer.Response[*serializer.ObjectProps]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`object/property/`)
		b.WriteString(q.ID)
		b.WriteString(`?trace_root=`)
		b.WriteString(strconv.FormatBool(q.TraceRoot))
		b.WriteString(`&is_folder=`)
		b.WriteString(strconv.FormatBool(q.IsFolder))
	}, nil, &out)
	return out.Data, err
}

// Share 分享

// 创建新分享 [分享参数] [分享链接,错误]
func (c *SiteObj) ShareNew(s *share.ShareCreateService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`share`) }, s, &out)
	return out.Data, err
}

// 列出我的分享 [查询参数] [分享列表,错误]
func (c *SiteObj) ShareGet(q *share.ShareListService) (*serializer.ShareList, error) {
	var out serializer.Response[*serializer.ShareList]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`share`) }, nil, &out)
	return out.Data, err
}

// 转存他人分享 (Pro) [分享id,转存参数] [错误]
func (c *SiteObj) ShareSave(id string, s *share.Service) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`share/save/`)
		b.WriteString(id)
	}, s, &out)
	return err
}

// 更新分享属性 [分享id,修改属性(password|preview_enabled)] [结果,错误]
func (c *SiteObj) ShareUpd(id string, s *share.ShareUpdateService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPatch, func(b *strings.Builder) {
		b.WriteString(`share/`)
		b.WriteString(id)
	}, s, &out)
	return out.Data, err
}

// 删除分享 [分享id] [错误]
func (c *SiteObj) ShareDel(id string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`share/`)
		b.WriteString(id)
	}, nil, &out)
	return err
}

// Tag 用户标签

// 创建文件分类标签 [标签参数] [标签id,错误]
func (c *SiteObj) TagFilter(s *explorer.FilterTagCreateService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`tag/filter`) }, s, &out)
	return out.Data, err
}

// 创建目录快捷方式标签 [标签参数] [标签id,错误]
func (c *SiteObj) TagLink(s *explorer.LinkTagCreateService) (*string, error) {
	var out serializer.Response[*string]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`tag/link`) }, s, &out)
	return out.Data, err
}

// 删除标签 [标签id] [错误]
func (c *SiteObj) TagDel(id string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`tag/`)
		b.WriteString(id)
	}, nil, &out)
	return err
}

// Vas 增值服务相关 (Pro)

// 获取容量包及配额信息 [] [信息,错误]
func (c *SiteObj) VasPack() (*serializer.Quota, error) {
	var out serializer.Response[*serializer.Quota]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`vas/pack`) }, nil, &out)
	return out.Data, err
}

// 获取商品信息，同时返回支付信息 [] [信息,错误]
func (c *SiteObj) VasProduct() (*serializer.ProductData, error) {
	var out serializer.Response[*serializer.ProductData]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`vas/product`) }, nil, &out)
	return out.Data, err
}

// 新建支付订单 [订单参数] [创建结果,错误]
func (c *SiteObj) VasOrderNew(s *vas.CreateOrderService) (*payment.OrderCreateRes, error) {
	var out serializer.Response[*payment.OrderCreateRes]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`vas/order`) }, s, &out)
	return out.Data, err
}

// 查询订单状态 [订单id] [状态(models.Order*),错误]
func (c *SiteObj) VasOrderGet(id string) (*int, error) {
	var out serializer.Response[*int]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`vas/order/`)
		b.WriteString(id)
	}, nil, &out)
	return out.Data, err
}

// 获取兑换码信息 [兑换码] [信息,错误]
func (c *SiteObj) VasRedeemGet(code string) (*vas.RedeemData, error) {
	var out serializer.Response[*vas.RedeemData]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) {
		b.WriteString(`vas/redeem/`)
		b.WriteString(code)
	}, nil, &out)
	return out.Data, err
}

// 执行兑换 [兑换码] [错误]
func (c *SiteObj) VasRedeemDo(code string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) {
		b.WriteString(`vas/redeem/`)
		b.WriteString(code)
	}, nil, &out)
	return err
}

// WebDAV 管理相关

// 获取账号信息 [] [账号列表,错误]
func (c *SiteObj) WebDavAccountsGet() (*setting.WebDAVAccountList, error) {
	var out serializer.Response[*setting.WebDAVAccountList]
	err := c.fetch(http.MethodGet, func(b *strings.Builder) { b.WriteString(`webdav/accounts`) }, nil, &out)
	return out.Data, err
}

// 新建账号 [账号参数] [账户信息,错误]
func (c *SiteObj) WebDavAccountsNew(s *setting.WebDAVAccountCreateService) (*setting.WebDAVAccountData, error) {
	var out serializer.Response[*setting.WebDAVAccountData]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`webdav/accounts`) }, s, &out)
	return out.Data, err
}

// 删除账号 [账号id] [错误]
func (c *SiteObj) WebDavAccountsDel(id string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`webdav/accounts`)
		b.WriteByte('/')
		b.WriteString(id)
	}, nil, &out)
	return err
}

// 删除目录挂载 [挂载id] [错误]
func (c *SiteObj) WebDavMountDel(id string) error {
	var out serializer.Response[struct{}]
	err := c.fetch(http.MethodDelete, func(b *strings.Builder) {
		b.WriteString(`webdav/mount`)
		b.WriteByte('/')
		b.WriteString(id)
	}, nil, &out)
	return err
}

// 创建目录挂载 [挂载参数] [挂载信息,错误]
func (c *SiteObj) WebDavMountNew(s *setting.WebDAVMountCreateService) (*setting.WebDAVAccountService, error) {
	var out serializer.Response[*setting.WebDAVAccountService]
	err := c.fetch(http.MethodPost, func(b *strings.Builder) { b.WriteString(`webdav/mount`) }, s, &out)
	return out.Data, err
}

// 更新账号可读性和是否使用代理服务 [更新参数] [结果(id×),错误]
func (c *SiteObj) WebDavUpd(s *setting.WebDAVAccountUpdateService) (*setting.WebDAVAccountUpdateService, error) {
	var out serializer.Response[*setting.WebDAVAccountUpdateService]
	err := c.fetch(http.MethodPatch, func(b *strings.Builder) { b.WriteString(`webdav/accounts`) }, s, &out)
	return out.Data, err
}
