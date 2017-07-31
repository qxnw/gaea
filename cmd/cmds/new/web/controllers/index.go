package controllers

var IndexControllerTmpl = `package controllers

import (
	"{@pImportName}/context"

	"github.com/qxnw/goplugin"
)

//IndexController 首页
type IndexController struct {
	requiredFields map[string][]string
}

//NewIndexController 创建首页
func NewIndexController() *IndexController {
	return &IndexController{}
}

//Handle 首页处理
func (c *IndexController) Handle(service string, ctx goplugin.Context, rpc goplugin.RPCInvoker) (status int, result interface{}, p map[string]interface{}, err error) {
	context, status, p, err := context.GetContext(ctx, rpc, c.requiredFields)
	if err != nil {
		return
	}
	defer context.Close()
	status = 200
	context.Info("--------访问首页--------")
	rt := make(map[string]string)
	err = context.IsLogin()
	if err == nil {
		context.ClearAuthTimes() //授权成功，清除授权次数
		rt["xsrf_token"] = context.GetXSRFToken()
		rt["xsrf_key"] = context.Conf.XSRFKey
		result = rt
		return
	}
	err = context.AddAuthTimes() //累加授权次数
	if err != nil {
		context.Errorf("授权超过限制次数：%v", err)
		context.Redirect(302, "/error")
		return
	}

	context.Error("用户登录信息不存在，跳转到微信授权")
	context.RedirectToWXAuth() //转跳到授权页面
	return
}

//Close 关闭服务
func (c *IndexController) Close() error {
	return nil
}
`
