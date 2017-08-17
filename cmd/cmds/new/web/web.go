package web

import (
	"github.com/qxnw/gaea/cmd/cmds/new/web/conf"
	"github.com/qxnw/gaea/cmd/cmds/new/web/controllers"
	"github.com/qxnw/gaea/cmd/cmds/new/web/views"
)

var mainTmpl = `package main

//go build -buildmode=plugin
import (
	"fmt"
	"github.com/qxnw/hydra/context"
)

//WebService 服务名称
type WebService struct {
}

//GetServices 获取当前插件提供的所有服务
func (p *WebService) GetServices() []string {
	return GetServices()
}

//Handle 业务处理
func (p *WebService) Handle(name string, mode string, service string, ctx *context.Context) (response context.Response, err error) {
	response, err = Handle(name, mode, service, ctx)
	if err != nil {
		err = fmt.Errorf("{@pImportName}:status:%d,err:%v", response.GetStatus(err), err)
	}
	return
}
func (p *WebService) Close() error {
	return nil
}

//GetWorker 获取当前worker
func GetWorker() context.Worker {
	return &WebService{}
}`

var TmplMap map[string]string

func init() {
	TmplMap = map[string]string{
		"main.go":              mainTmpl,
		"routers.go":           routerTmpl,
		"views/index.html":     views.IndexViewTmpl,
		"controllers/index.go": controllers.IndexControllerTmpl,
		"conf/conf.go":         conf.ConfTmpl,
	}
}
