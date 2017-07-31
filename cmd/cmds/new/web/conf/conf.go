package conf

var ConfTmpl = `package conf

import "github.com/qxnw/goplugin"
import "encoding/json"
import "github.com/qxnw/lib4go/concurrent/cmap"
import "fmt"

type {@pClassName}Conf struct {
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
func GetConf(ctx *goplugin.PluginContext) (c *{@pClassName}Conf, err error) {
	name, err := ctx.GetArgByName("conf")
	if err != nil {
		return nil, err
	}
	_, v, err := confCache.SetIfAbsentCb(name, func(input ...interface{}) (interface{}, error) {
		content, err := ctx.GetVarParam("conf", name)
		if err != nil {
			return nil, err
		}
		conf := &{@pClassName}Conf{}
		err = json.Unmarshal([]byte(content), conf)
		if err != nil {
			err = fmt.Errorf("conf配置文件错误:%v", err)
			return nil, err
		}
		return conf, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*{@pClassName}Conf), nil
}
`
