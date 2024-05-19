package cr

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ZxwyWebSite/cr-go-sdk/pkg/json"
	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
)

type UploadTask struct {
	Site    *SiteObj
	File    io.Reader
	Size    uint64 // 文件大小
	Name    string // 文件名称 (指定时为上传名称,否则使用原名)
	Mime    string // 文件类型 (先通过读文件获取,再判断扩展名)
	Sess    *serializer.UploadCredential
	Policy  *serializer.PolicySummary
	ModTime int64 // 创建时间 (UnixMilli)
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
			c.Sess.SessionID, chunk, &io.LimitedReader{R: c.File, N: int64(byteSize)}, byteSize, c.Mime,
		)
		if err == nil {
			try = 0
			chunk++
			finish += byteSize
		}
	}
	return err
}

// 从机存储
func (c *UploadTask) remote() error {
	uploadUrl := c.Sess.UploadURLs[0]
	var chunk int
	var oerr error
	for try := 0; try < 3; try++ {
		var n int64 = int64(c.Size) - int64(c.Sess.ChunkSize)*int64(chunk)
		if n <= 0 {
			break
		}
		if n > int64(c.Sess.ChunkSize) {
			n = int64(c.Sess.ChunkSize)
		}
		var b strings.Builder
		b.WriteString(uploadUrl)
		b.WriteString(`?chunk=`)
		b.WriteString(strconv.Itoa(chunk))
		req, err := http.NewRequest(http.MethodPost, b.String(), io.LimitReader(c.File, n))
		if err != nil {
			return err
		}
		req.Header.Set(`Authorization`, c.Sess.Credential)
		req.ContentLength = n
		// req.Header.Set(`Content-Length`, strconv.FormatUint(n, 10))
		req.Header.Set(`Content-Type`, `application/octet-stream`)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			oerr = err
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
		chunk++
	}
	return oerr
}

// OneDrive
func (c *UploadTask) onedrive() error {
	uploadUrl := c.Sess.UploadURLs[0]
	var finish uint64 = 0
	for finish < c.Size {
		var byteSize = c.Sess.ChunkSize
		left := c.Size - finish
		if left < c.Sess.ChunkSize {
			byteSize = left
		}
		req, err := http.NewRequest(http.MethodPut, uploadUrl, io.LimitReader(c.File, int64(byteSize)))
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
		req.ContentLength = int64(byteSize)
		// req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
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
		req, err := http.NewRequest(http.MethodPut, uploadUrl, io.LimitReader(c.File, int64(byteSize)))
		if err != nil {
			return err
		}
		req.ContentLength = int64(byteSize)
		// req.Header.Set(`Content-Length`, strconv.FormatUint(byteSize, 10))
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
	req, err := http.NewRequest(http.MethodPost, c.Sess.CompleteURL, strings.NewReader(b.String()))
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
	return c.Site.CallbackS3(c.Sess.SessionID)
}

// 执行上传 [上传目录] [错误]
func (c *UploadTask) Do(dir string) error {
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
			println(string(s))
		}
		sess, err := c.Site.FileUploadNew(sreq)
		if err != nil {
			return err
		}
		if Cr_Debug {
			s, _ := json.MarshalIndent(sess, ``, `    `)
			println(string(s))
		}
		if sess.ChunkSize == 0 {
			sess.ChunkSize = c.Size
		}
		c.Sess = sess
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
