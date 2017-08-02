package controllers

var IndexControllerTmpl = `package controllers

import (
	"{@pImportName}/context"

	"github.com/qxnw/goplugin"
)

//IndexController 首页
type IndexController struct {
	fields map[string][]string
}

//NewIndexController 创建首页
func NewIndexController() *IndexController {
	return &IndexController{}
}

//Handle 首页处理
func (c *IndexController) Handle(service string, ctx goplugin.Context, rpc goplugin.RPCInvoker) (status int, result interface{}, p map[string]interface{}, err error) {
	context, status, p, err := context.GetContext(ctx, rpc, c.fields)
	if err != nil {
		return
	}
	defer context.Close()
	status = 200
	context.Info("--------访问首页--------")
	rt := make(map[string]string)
	err = context.IsLogin()
	if err == nil {	
		rt["xsrf_token"] = context.GetXSRFToken()
		rt["xsrf_key"] = context.Conf.XSRFKey
		result = rt
		return
	}
	return
}

//Close 关闭服务
func (c *IndexController) Close() error {
	return nil
}
`
