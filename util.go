package cr

import (
	"fmt"
	"math"
	"net/http"
	_ "unsafe"

	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
)

/*func jsonReader(obj any) (io.Reader, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	// if err != nil {
	// 	return nil, err
	// }
	return &buf, err
}*/

// 格式化文件大小
func SizeToString(bytes uint64) string {
	if bytes == 0 {
		return "0B"
	}
	const k = 1024
	sizes := []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))
	return fmt.Sprintf("%.1f%s", float64(bytes)/math.Pow(k, float64(i)), sizes[i])
}

/*func SizeToString(bytes uint64) string {
	if bytes == 0 {
		return "0 B"
	}
	const k = 1024
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))

	return fmt.Sprintf("%.1f %s", float64(bytes)/math.Pow(k, float64(i)), sizes[i])
}*/

//go:linkname readSetCookies net/http.readSetCookies
func readSetCookies(h http.Header) []*http.Cookie

/*func ParseCookies(strs ...string) []*http.Cookie {
	return readSetCookies(http.Header{
		`Set-Cookie`: strs,
	})
}*/

// 格式化Cookie字符串
func ParseCookie(str string) *http.Cookie {
	return readSetCookies(http.Header{
		`Set-Cookie`: []string{str},
	})[0]
}

// 生成文件列表
func GenerateSrc(isDir bool, ids ...string) *explorer.ItemIDService {
	var list explorer.ItemIDService
	if isDir {
		list.Dirs = ids
	} else {
		list.Items = ids
	}
	return &list
}
