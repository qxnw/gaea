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
		response.SetContent(context.ERR_NOT_EXTENDED, err)
		return
	}
	status, err := ctx.Input.Check(c.fields)
	if err != nil {
		response.SetContent(status, err)
		return
	}
	ctx.Info("----------访问首页----------")
	ctx.Info(config)
	return
}

//Close 关闭服务
func (c *IndexController) Close() error {
	return nil
}`
