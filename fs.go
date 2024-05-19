package cr

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"time"

	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
)

type FSFileInfo struct {
	obj  *serializer.Object
	prop *serializer.ObjectProps
}

func (f *FSFileInfo) IsDir() bool {
	return f.obj.Type == `dir`
}

func (f *FSFileInfo) ModTime() time.Time {
	return f.prop.CreatedAt
}

func (f *FSFileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

func (f *FSFileInfo) Name() string {
	return f.obj.Name
}

func (f *FSFileInfo) Size() int64 {
	return int64(f.prop.Size)
}

func (f *FSFileInfo) Sys() any {
	return nil
}

type FSFile struct {
	site *SiteObj
	obj  *serializer.Object
	file io.ReadCloser
}

func (f *FSFile) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FSFile) Read(p []byte) (int, error) {
	if f.file == nil {
		link, err := f.site.FileDownload(f.obj.ID)
		if err != nil {
			return 0, err
		}
		if (*link)[0] == '/' {
			*link = f.site.Addr + (*link)[1:]
		}
		req, err := http.NewRequest(http.MethodGet, *link, nil)
		if err != nil {
			return 0, err
		}
		req.Header[`User-Agent`] = []string{Cr_UserAgent}
		req.Header[`Accept`] = []string{Cr_Accept}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
		f.file = res.Body
	}
	return f.file.Read(p)
}

func (f *FSFile) Stat() (fs.FileInfo, error) {
	prop, err := f.site.ObjectProperty(&explorer.ItemPropertyService{ID: f.obj.ID})
	if err != nil {
		return nil, err
	}
	return &FSFileInfo{f.obj, prop}, nil
}

type FS SiteObj

func (c FS) Open(name string) (fs.File, error) {
	site := (*SiteObj)(&c)
	list, err := site.Directory(path.Dir(name))
	if err != nil {
		return nil, err
	}
	name = path.Base(name)
	var obj *serializer.Object
	for _, v := range list.Objects {
		if v.Name == name && v.Type == `file` {
			obj = &v
			break
		}
	}
	if obj == nil {
		return nil, fs.ErrNotExist
	}
	return &FSFile{site, obj, nil}, nil
}
