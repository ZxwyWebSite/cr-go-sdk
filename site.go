package cr

import (
	"strings"

	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
)

// Api版本 访问 站点地址/api/v3/site/ping 获取
const (
	ApiPlus uint8 = iota // 3.8.3+1.1-plus
	ApiV353              // 3.5.3
	ApiV383              // 3.8.3
	// ApiV400 // 4.0.0
)

// Cloudreve 站点驱动
type SiteObj struct {
	// Client  *http.Client // HTTP 客户端
	Addr    string                 // 站点地址 (首页)
	ApiVer  uint8                  // 接口版本
	Config  *serializer.SiteConfig // 站点配置
	Users   *UserObj               // 账号数据
	Version string                 // 程序版本
}

// 新建站点对象
func NewSite(addr string, apiver uint8) (*SiteObj, error) {
	if !strings.HasSuffix(addr, `/`) {
		addr += `/`
	}
	site := &SiteObj{
		Addr:   addr,
		ApiVer: apiver,
	}
	/*config, err := site.SiteConfig()
	if err == nil && config != nil {
		site.Config = config
	}
	version, err := site.SitePing()
	if err == nil && version != nil {
		site.Version = *version
	}*/
	return site, nil
}

/*func (c *SiteObj) Login(mail, pass string) (*UserObj, error) {
	user := UserObj{
		Mail: mail,
		Pass: pass,
	}
	c.Users = append(c.Users, user)
	return &user, nil
}*/
