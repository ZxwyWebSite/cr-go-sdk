/*
A Simple CloudreveV3 SDK
*/
package cr

import (
	"fmt"
	"net/http"
	"os"
)

// 一些全局常量
const (
	Cr_Version   = `0.0.3`
	Cr_UserAgent = `Mozilla/5.0 (compatible; cr-go-sdk/` + Cr_Version + `)`
	Cr_Accept    = `application/json, text/plain, */*`
)

// 一些全局变量
var (
	Cr_Debug = true

	Cr_OcrApi   = `https://api.nn.ci/ocr/file/json`
	Cr_OcrRetry = 10

	Cr_Client = http.DefaultClient
	Cr_Format = func(args ...any) { fmt.Fprintln(os.Stderr, args...) }
)
