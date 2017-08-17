package controllers

var IndexControllerTmpl = `package controllers

import (
	"{@pImportName}/conf"

	"github.com/qxnw/hydra/context"
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
func (c *IndexController) Handle(name string, mode string, service string, ctx *context.Context) (response *context.WebResponse, err error) {
	response = context.GetWebResponse(ctx)
	config, err := conf.GetConf(ctx)
	if err != nil {
		response.SetError(context.ERR_NOT_EXTENDED, err)
		return
	}
	status, err := ctx.Input.Check(c.fields)
	if err != nil {
		response.SetError(status, err)
		return
	}
	ctx.Info("----------访问首页----------")
	rt := make(map[string]string)
	err = response.IsLogin()
	if err == nil {
		response.ClearAuthTimes() //授权成功，清除授权次数
		rt["xsrf_token"] = response.MakeXSRFToken(config.XSRFSecret)
		rt["xsrf_key"] = config.XSRFKey
		response.Success(rt)
		return
	}
	err = response.AddAuthTimes() //累加授权次数
	if err != nil {
		ctx.Errorf("授权超过限制次数：%v", err)
		response.Redirect(302, "/error")
		return
	}

	ctx.Error("用户登录信息不存在，跳转到微信授权")
	response.RedirectToWXAuth(config.WXAuthNotifyURL, config.APPID, config.WXAuthNotifyURL) //转跳到授权页面
	return
}

//Close 关闭服务
func (c *IndexController) Close() error {
	return nil
}`
