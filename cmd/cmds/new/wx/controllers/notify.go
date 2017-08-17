package controllers

var NotifyControllerTmpl = `package controllers

import (
	"{@pImportName}/conf"

	"github.com/qxnw/hydra/context"
)

//NotifyController 接收微信通知
type NotifyController struct {
	fields map[string][]string
}

//NewIndexController 创建首页
func NewNotifyController() *NotifyController {
	return &NotifyController{
		fields: map[string][]string{
			"input": []string{"code"},
		},
	}
}

//Handle 首页处理
func (c *NotifyController) Handle(name string, mode string, service string, ctx *context.Context) (response *context.WebResponse, err error) {
	response = context.GetWebResponse(ctx)
	config, err := conf.GetConf(ctx)
	if err != nil {
		response.SetError(context.ERR_NOT_EXTENDED, err)
		return
	}
	//检查输入参数
	status, err := ctx.Input.Check(c.fields)
	if err != nil {
		response.SetError(status, err)
		return
	}

	//使用code登录
	err = response.Login(config.AuthorizeLogin)
	if err != nil {
		return
	}
	ctx.Error("用户登录信息不存在，跳转到微信授权")
	response.RedirectToWXAuth(config.WXAuthNotifyURL, config.APPID, config.WXAuthNotifyURL) //转跳到授权页面
	return
}

//Close 关闭服务
func (c *NotifyController) Close() error {
	return nil
}`
