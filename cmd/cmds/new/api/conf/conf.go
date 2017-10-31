package conf

var ConfTmpl = `import (
	"encoding/json"
	"fmt"

	"github.com/qxnw/hydra/context"

	"github.com/qxnw/lib4go/concurrent/cmap"
)

type APIConf struct {
}

var confCache cmap.ConcurrentMap

func init() {
	confCache = cmap.New(3)
}

//GetConf 获取系统配置文件
func GetConf(ctx *context.Context) (c *APIConf, err error) {
	name, err := ctx.Input.GetArgsByName("conf")
	if err != nil {
		return nil, err
	}
	_, v, err := confCache.SetIfAbsentCb(name, func(input ...interface{}) (interface{}, error) {
		name := input[0].(string)
		content, err := ctx.Input.GetVarParam("conf", name)
		if err != nil {
			return nil, err
		}
		conf := &APIConf{}
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
	return v.(*APIConf), nil
}`
