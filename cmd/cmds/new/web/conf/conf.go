package conf

var ConfTmpl = `package conf

import (
	"encoding/json"
	"fmt"

	"github.com/qxnw/hydra/context"

	"github.com/qxnw/lib4go/concurrent/cmap"
)

type WebConf struct {
	AuthorizeLogin  string
	TaskShareURL    string
	WXAuthURL       string
	APPID           string
	HostName        string
	WXAuthNotifyURL string
	XSRFKey         string
	XSRFSecret      string
}

var confCache cmap.ConcurrentMap

func init() {
	confCache = cmap.New(3)
}

//GetConf 获取小微配置信息
func GetConf(ctx *context.Context) (c *WebConf, err error) {
	name, err := ctx.Request.Setting.Get("conf")
	if err != nil {
		return nil, err
	}
	_, v, err := confCache.SetIfAbsentCb(name, func(input ...interface{}) (interface{}, error) {
		name := input[0].(string)
		content, err := ctx.Input.GetVarParam("conf", name)
		if err != nil {
			return nil, err
		}
		conf := &WebConf{}
		err = json.Unmarshal([]byte(content), conf)
		if err != nil {
			err = fmt.Errorf("conf配置文件错误:%v", err)
			return nil, err
		}
		return conf, nil
	}, name)
	if err != nil {
		return nil, err
	}
	return v.(*WebConf), nil
}`
