package cr

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	json "github.com/ZxwyProject/zson"
	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/user"
)

// (SDK) 初始化站点数据 [] [错误]
func (c *SiteObj) SdkInit() error {
	if Cr_Debug {
		Cr_Format(`[sdk-debug] Cr_Debug is enabled, set cr.Cr_Debug=false to disable it`)
	}
	// 注：调用SiteConfig似乎可以刷新Cookie过期时间
	config, err := c.SiteConfig()
	if err == nil {
		if config != nil {
			c.Config = config
		}
	} else {
		return err
	}
	version, err := c.SitePing()
	if err == nil {
		if version != nil {
			c.Version = *version
		}
	} else {
		return err
	}
	/*if c.Users != nil && !c.Config.LoginCaptcha {
		if c.Users.Sess == `` || c.Config.User.Anonymous {
			user, err := c.UserSession(&user.LoginInfo{
				UserName: c.Users.Mail,
				Password: c.Users.Pass,
				// UserLoginService: user.UserLoginService{
				// 	UserName: c.Users.Mail,
				// 	Password: c.Users.Pass,
				// },
			})
			if err == nil {
				if user != nil {
					c.Config.User = *user
					c.Users.Exps = time.Now().Add(time.Hour * 24 * 7).Unix()
					println(`[debug] userlogin:`, c.Users.Sess)
				}
			} else {
				return err
			}
		}
	}*/
	return nil
}

// (SDK) 上传文件 [上传目录,文件,上传名称] [错误]
func (c *SiteObj) SdkUpload(dir string, file *os.File, name string) error {
	// 获取必要信息
	/*
		info: 文件信息 fs.FileInfo
		size: 文件大小 uint64
		name: 文件名称 string (指定时为上传名称,否则使用原名)
		mime: 文件类型 string (先通过读文件获取,再判断扩展名)
		list: 目录信息 serializer.ObjectList (查找存储策略)
		sess: 上传会话 serializer.UploadCredential
	*/
	info, err := file.Stat()
	if err != nil {
		return err
	}
	size := uint64(info.Size())
	if name == `` {
		name = info.Name()
	}
	mbuf := make([]byte, 512)
	mn, err := file.Read(mbuf)
	if err != nil {
		return err
	}
	mime := http.DetectContentType(mbuf[:mn])
	/*if mime == `application/octet-stream` {
		if me := stdmime.TypeByExtension(filepath.Ext(name)); me != `` {
			mime = me
		}
	}*/
	_, err = file.Seek(0, io.SeekStart) // 重置reader指针!
	if err != nil {
		return err
	}
	task := &UploadTask{
		Site:    c,
		File:    file,
		Size:    size,
		Name:    name,
		Mime:    mime,
		ModTime: info.ModTime().UnixMilli(),
	}
	if Cr_Debug {
		s, _ := json.MarshalIndent(task, ``, `    `)
		Cr_Format(string(s))
	}
	return task.Do(dir)
	/*list, err := c.Directory(dir)
	if err != nil {
		return err
	}
	sreq := &explorer.CreateUploadSessionService{
		Path:         dir,
		Size:         size,
		Name:         name,
		PolicyID:     list.Policy.ID,
		LastModified: info.ModTime().UnixMilli(),
		MimeType:     mime,
	}
	if Cr_Debug {
		s, _ := json.MarshalIndent(sreq, ``, `    `)
		println(string(s))
	}
	sess, err := c.FileUploadNew(sreq)
	if err != nil {
		return err
	}
	if Cr_Debug {
		s, _ := json.MarshalIndent(sess, ``, `    `)
		println(string(s))
	}
	switch list.Policy.Type {
	case `local`: // 本机存储 参考alist/cloudreve驱动
		var buf []byte
		var chunk int
		for {
			var n int
			buf = make([]byte, sess.ChunkSize)
			n, err = io.ReadAtLeast(file, buf, int(sess.ChunkSize))
			if err != nil && err != io.ErrUnexpectedEOF {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if n == 0 {
				break
			}
			err = c.FileUploadPut(sess.SessionID, strconv.Itoa(chunk), buf[:n], mime)
			if err != nil {
				break
			}
			chunk++
		}
	case `remote`: // 从机存储 自行抓包分析
		uploadUrl := sess.UploadURLs[0]
		var chunk int
		for {
			var n int
			buf := make([]byte, sess.ChunkSize)
			n, err = io.ReadAtLeast(file, buf, int(sess.ChunkSize))
			if err != nil && err != io.ErrUnexpectedEOF {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if n == 0 {
				break
			}
			buf = buf[:n]
			var b strings.Builder
			b.WriteString(uploadUrl)
			b.WriteString(`?chunk=`)
			b.WriteString(strconv.Itoa(chunk))
			req, err := http.NewRequest(http.MethodPost, b.String(), bytes.NewBuffer(buf))
			if err != nil {
				return err
			}
			req.Header.Set(`Authorization`, sess.Credential)
			req.Header.Set(`Content-Length`, strconv.Itoa(n))
			req.Header.Set(`Content-Type`, `application/octet-stream`)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			var out serializer.Response[struct{}]
			if err = json.NewDecoder(res.Body).Decode(&out); err == nil {
				err = out.Err()
			}
			if err != nil {
				res.Body.Close()
				return err
			}
			res.Body.Close()
			chunk++
		}
	case `onedrive`: // OneDrive 参考alist/onedrive驱动
		uploadUrl := sess.UploadURLs[0]
		var finish uint64 = 0
		for finish < size {
			var byteSize = sess.ChunkSize
			left := size - finish
			if left < sess.ChunkSize {
				byteSize = left
			}
			byteData := make([]byte, byteSize)
			_, err := io.ReadFull(file, byteData)
			if err != nil {
				return err
			}
			req, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewBuffer(byteData))
			if err != nil {
				return err
			}
			var b strings.Builder
			b.WriteString(`bytes `)
			b.WriteString(strconv.FormatUint(finish, 10))
			b.WriteByte('-')
			finish += byteSize
			b.WriteString(strconv.FormatUint(finish-1, 10))
			b.WriteByte('/')
			b.WriteString(strconv.FormatUint(size, 10))
			req.Header.Set(`Content-Range`, b.String())
			req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
			req.Header.Set(`Content-Type`, `application/octet-stream`)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if res.StatusCode != 201 && res.StatusCode != 202 && res.StatusCode != 200 {
				data, _ := io.ReadAll(res.Body)
				res.Body.Close()
				return errors.New(string(data))
			}
			res.Body.Close()
		}
		return c.CallbackOneDriveFinish(sess.SessionID)
	case `s3`: // AWS S3 自行抓包分析
		var finish uint64 = 0
		var etag = make([]string, len(sess.UploadURLs))
		for chunk, uploadUrl := range sess.UploadURLs {
			var byteSize = sess.ChunkSize
			left := size - finish
			if left < sess.ChunkSize {
				byteSize = left
			}
			byteData := make([]byte, byteSize)
			n, err := io.ReadFull(file, byteData)
			if err != nil && err != io.ErrUnexpectedEOF {
				if err == io.EOF {
					break
				}
				return err
			}
			req, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewBuffer(byteData[:n]))
			if err != nil {
				return err
			}
			req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
			req.Header.Set(`Content-Type`, `application/octet-stream`)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			etag[chunk] = res.Header.Get(`Etag`)
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}
		var b strings.Builder
		b.WriteString(`<CompleteMultipartUpload>`)
		for i, e := range etag {
			b.WriteString(`<Part><PartNumber>`)
			b.WriteString(strconv.Itoa(i + 1))
			b.WriteString(`</PartNumber><ETag>`)
			b.WriteString(e)
			b.WriteString(`</ETag></Part>`)
		}
		b.WriteString(`</CompleteMultipartUpload>`)
		req, err := http.NewRequest(http.MethodPost, sess.CompleteURL, strings.NewReader(b.String()))
		if err != nil {
			return err
		}
		req.Header.Set(`Content-Type`, `application/xhtml+xml`)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			data, _ := io.ReadAll(res.Body)
			res.Body.Close()
			return errors.New(string(data))
		}
		res.Body.Close()
		return c.CallbackS3(sess.SessionID)
	default:
		return errors.New(`不支持的存储策略类型`)
	}
	return nil*/
}

// (SDK) 识别验证码 (仅支持默认版) [] [结果,错误]
func (c *SiteObj) SdkSolveCaptcha() (*string, error) {
	captcha, err := c.SiteCaptcha()
	if err != nil {
		return nil, err
	}
	decs, err := base64.StdEncoding.DecodeString((*captcha)[strings.Index(*captcha, `,`)+1:])
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	var b strings.Builder
	b.WriteString(`form-data; name="`)
	b.WriteString(`image`)
	b.WriteString(`"; filename="`)
	b.WriteString(`validateCode.png`)
	b.WriteByte('"')
	h.Set(`Content-Disposition`, b.String())
	h.Set(`Content-Type`, `image/png`)
	pw, err := mw.CreatePart(h)
	if err == nil {
		_, err = pw.Write(decs)
		if err == nil {
			err = mw.Close()
		}
	}
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(http.MethodPost, Cr_OcrApi, buf, false)
	if err != nil {
		return nil, err
	}
	req.Header[`Content-Type`] = []string{mw.FormDataContentType()}
	res, err := Cr_Client.Do(req)
	if err != nil {
		return nil, err
	}
	var out struct {
		Status int    `json:"status"`
		Result string `json:"result"`
		Msg    string `json:"msg"`
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		res.Body.Close()
		return nil, err
	}
	res.Body.Close()
	if out.Status != 200 {
		return nil, errors.New(`ocr: ` + out.Msg)
	}
	return &out.Result, nil
}

// (SDK) 登录账号 [] [错误]
func (c *SiteObj) SdkLogin() error {
	if c.Users == nil {
		return errors.New(`未添加账号信息`)
	}
	var req = &user.LoginInfo{
		UserName: c.Users.Mail,
		Password: c.Users.Pass,
	}
	if c.Config.LoginCaptcha {
		if c.Config.CaptchaType != `normal` {
			return errors.New(`暂不支持该类型验证码: ` + c.Config.CaptchaType)
		}
		var err error
		for i := 0; i < Cr_OcrRetry; i++ {
			var result *string
			result, err = c.SdkSolveCaptcha()
			if err != nil {
				continue
			}
			req.CaptchaCode = *result
			var user *serializer.User
			user, err = c.UserSession(req)
			if err != nil {
				if apperr, ok := err.(*serializer.AppError); ok {
					if apperr.Code == serializer.CodeCaptchaError || apperr.Msg == `验证码错误` {
						continue
					}
				}
				return err
			}
			c.Config.User = *user
			break
		}
		return err
	}
	user, err := c.UserSession(req)
	if err == nil {
		c.Config.User = *user
	}
	return err
}

// (SDK) 登录带验证码的站点
/*func (c *SiteObj) SdkLoginWithCaptcha() error {
	if !c.Config.LoginCaptcha {
		return errors.New(`该站点未启用验证码`)
	}
	if c.Config.CaptchaType != `normal` {
		return errors.New(`暂不支持该类型验证码: ` + c.Config.CaptchaType)
	}
	for {
		captcha, err := c.SiteCaptcha()
		if err != nil {
			return err
		}
		if captcha == nil {
			return errors.New(`can not get captcha`)
		}
		fmt.Println(*captcha)
		dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader((*captcha)[strings.Index(*captcha, `,`)+1:]))
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		h := make(textproto.MIMEHeader)
		var b strings.Builder
		b.WriteString(`form-data; name="`)
		b.WriteString(`image`)
		b.WriteString(`"; filename="`)
		b.WriteString(`validateCode.png`)
		b.WriteByte('"')
		h.Set(`Content-Disposition`, b.String())
		h.Set(`Content-Type`, `image/png`)
		pw, err := mw.CreatePart(h)
		if err != nil {
			return err
		}
		_, err = io.Copy(pw, dec)
		if err != nil {
			return err
		}
		mw.Close()
		req, err := http.NewRequest(http.MethodPost, `https://api.nn.ci/ocr/file/json`, buf)
		if err != nil {
			return err
		}
		req.Header.Set(`Content-Type`, mw.FormDataContentType())
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		var out struct {
			Status int    `json:"status"`
			Result string `json:"result"`
			Msg    string `json:"msg"`
		}
		if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
			res.Body.Close()
			return err
		}
		res.Body.Close()
		if out.Status != 200 {
			return errors.New(`ocr error: ` + out.Msg)
		}
		fmt.Printf("%+v\n", out)
		if len(out.Result) != 6 {
			continue
		}
		user, err := c.UserSession(&user.LoginInfo{
			UserName:    c.Users.Mail,
			Password:    c.Users.Pass,
			CaptchaCode: out.Result,
		})
		if err != nil {
			if apperr, ok := err.(*serializer.AppError); ok {
				if apperr.Code == 40001 || apperr.Msg == `验证码错误` {
					fmt.Println(apperr.Error())
					continue
				}
			}
			return err
		}
		c.Config.User = *user
		return nil
	}
}*/
