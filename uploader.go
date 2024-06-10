package cr

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	json "github.com/ZxwyProject/zson"
	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
)

// io.LimitedReader with callback
type uploadReader struct {
	R io.Reader
	N uint64
	C func(n int, e error)
}

func (l *uploadReader) Read(p []byte) (n int, err error) {
	if l.N == 0 {
		l.C(0, io.EOF)
		return 0, io.EOF
	}
	if uint64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	// 防止 uint 被减到负数溢出?
	if uint64(n) >= l.N {
		l.N = 0
	} else {
		l.N -= uint64(n)
	}
	l.C(n, err)
	return
}

type UploadTask struct {
	Site    *SiteObj
	File    io.Reader
	Size    uint64 // 文件大小
	Name    string // 文件名称 (指定时为上传名称,否则使用原名)
	Mime    string // 文件类型 (先通过读文件获取,再判断扩展名)
	Sess    *serializer.UploadCredential
	Policy  *serializer.PolicySummary
	ModTime int64 // 创建时间 (UnixMilli)

	Chunks   []uint64                        // 预切片大小
	Callback func(chunk int, n int, e error) // 进度回调
}

// 本机存储
func (c *UploadTask) local() error {
	var chunk int
	var finish uint64 = 0
	var err error
	for try := 0; try < 3; try++ {
		if finish >= c.Size {
			break
		}
		var byteSize = c.Sess.ChunkSize
		left := c.Size - finish
		if left < c.Sess.ChunkSize {
			byteSize = left
		}
		err = c.Site.FileUploadPut(
			c.Sess.SessionID, chunk, &uploadReader{R: c.File, N: byteSize, C: func(n int, e error) {
				c.Callback(chunk, n, e)
			}}, byteSize, c.Mime,
		)
		if err == nil {
			try = 0
			chunk++
			finish += byteSize
		}
		if s, ok := c.File.(io.ReadSeeker); ok {
			_, err = s.Seek(int64(finish), io.SeekStart)
			if err != nil {
				break
			}
		} else {
			// 不重置Reader无法重试
			break
		}
	}
	return err
}

/*func (c *UploadTask) localsync() error {
	return nil
}*/

// 从机存储
func (c *UploadTask) remote() error {
	uploadUrl := c.Sess.UploadURLs[0]
	for chunk, size := range c.Chunks {
		var oerr error
		for try := 0; try < 3; try++ {
			var b strings.Builder
			b.WriteString(uploadUrl)
			b.WriteString(`?chunk=`)
			b.WriteString(strconv.Itoa(chunk))
			req, err := c.Site.newRequest(http.MethodPost, b.String(), &uploadReader{c.File, size, func(n int, e error) {
				c.Callback(chunk, n, e)
			}}, false)
			if err != nil {
				return err
			}
			req.Header.Set(`Authorization`, c.Sess.Credential)
			req.ContentLength = int64(size) // req.Header.Set(`Content-Length`, strconv.FormatUint(n, 10))
			req.Header.Set(`Content-Type`, `application/octet-stream`)
			res, err := Cr_Client.Do(req)
			if err != nil {
				if s, ok := c.File.(io.ReadSeeker); ok {
					_, oerr = s.Seek(int64(c.Sess.ChunkSize*uint64(chunk-1)+size), io.SeekStart)
					if oerr != nil {
						break
					}
				} else {
					// 不重置Reader无法重试
					return err
				}
				continue
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
			try = 0
			break
		}
		if oerr != nil {
			return oerr
		}
	}
	return nil
}

// OneDrive
func (c *UploadTask) onedrive() error {
	uploadUrl := c.Sess.UploadURLs[0]
	var chunk int
	var finish uint64 = 0
	for finish < c.Size {
		var byteSize = c.Sess.ChunkSize
		left := c.Size - finish
		if left < c.Sess.ChunkSize {
			byteSize = left
		}
		req, err := c.Site.newRequest(http.MethodPut, uploadUrl, &uploadReader{c.File, byteSize, func(n int, e error) {
			c.Callback(chunk, n, e)
		}}, false)
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
		b.WriteString(strconv.FormatUint(c.Size, 10))
		req.Header.Set(`Content-Range`, b.String())
		req.ContentLength = int64(byteSize) // req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
		req.Header.Set(`Content-Type`, `application/octet-stream`)
		res, err := Cr_Client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 201 && res.StatusCode != 202 && res.StatusCode != 200 {
			data, _ := io.ReadAll(res.Body)
			res.Body.Close()
			return errors.New(string(data))
		}
		res.Body.Close()
		chunk++
	}
	return c.Site.CallbackOneDriveFinish(c.Sess.SessionID)
}

// AWS S3
func (c *UploadTask) s3() error {
	var finish uint64 = 0
	var etag = make([]string, len(c.Sess.UploadURLs))
	for chunk, uploadUrl := range c.Sess.UploadURLs {
		var byteSize = c.Sess.ChunkSize
		left := c.Size - finish
		if left < c.Sess.ChunkSize {
			byteSize = left
		}
		req, err := c.Site.newRequest(http.MethodPut, uploadUrl, &uploadReader{c.File, byteSize, func(n int, e error) {
			c.Callback(chunk, n, e)
		}}, false)
		if err != nil {
			return err
		}
		req.ContentLength = int64(byteSize) // req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
		req.Header.Set(`Content-Type`, `application/octet-stream`)
		res, err := Cr_Client.Do(req)
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
	req, err := c.Site.newRequest(http.MethodPost, c.Sess.CompleteURL, strings.NewReader(b.String()), false)
	if err != nil {
		return err
	}
	req.Header.Set(`Content-Type`, `application/xhtml+xml`)
	res, err := Cr_Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		data, _ := io.ReadAll(res.Body)
		res.Body.Close()
		return errors.New(string(data))
	}
	res.Body.Close()
	return c.Site.CallbackS3(c.Sess.SessionID)
}

// 创建任务(补全参数) [上传目录] [错误]
func (c *UploadTask) In(dir string) error {
	if c.Policy == nil {
		list, err := c.Site.Directory(dir)
		if err != nil {
			return err
		}
		c.Policy = list.Policy
	}
	if c.ModTime == 0 {
		c.ModTime = time.Now().UnixMilli()
	}
	if c.Sess == nil {
		sreq := &explorer.CreateUploadSessionService{
			Path:         dir,
			Size:         c.Size,
			Name:         c.Name,
			PolicyID:     c.Policy.ID,
			LastModified: c.ModTime,
			MimeType:     c.Mime,
		}
		if Cr_Debug {
			s, _ := json.MarshalIndent(sreq, ``, `    `)
			Cr_Format(string(s))
		}
		sess, err := c.Site.FileUploadNew(sreq)
		if err != nil {
			return err
		}
		if Cr_Debug {
			s, _ := json.MarshalIndent(sess, ``, `    `)
			Cr_Format(string(s))
		}
		if sess.ChunkSize == 0 {
			sess.ChunkSize = c.Size
		}
		c.Sess = sess
	}
	for total := c.Size; true; {
		if total > c.Sess.ChunkSize {
			c.Chunks = append(c.Chunks, c.Sess.ChunkSize)
			total -= c.Sess.ChunkSize
		} else {
			c.Chunks = append(c.Chunks, total)
			break
		}
	}
	return nil
}
func (c *UploadTask) Go() error {
	if c.Callback == nil {
		c.Callback = func(chunk, n int, e error) {}
	}
	switch c.Policy.Type {
	case `local`:
		return c.local()
	case `remote`:
		return c.remote()
	case `onedrive`:
		return c.onedrive()
	case `s3`:
		return c.s3()
	default:
		return errors.New(`暂不支持当前存储策略类型: ` + c.Policy.Type)
	}
}

// 执行上传 [上传目录] [错误]
func (c *UploadTask) Do(dir string) error {
	if err := c.In(dir); err != nil {
		return err
	}
	return c.Go()
}
