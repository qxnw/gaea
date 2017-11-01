package api

import (
	"github.com/qxnw/gaea/cmd/cmds/new/api/conf"
	"github.com/qxnw/gaea/cmd/cmds/new/api/services"
)

var mainTmpl = `package main

//go build -buildmode=plugin
import (
	"fmt"
	"github.com/qxnw/hydra/context"
)

//{@pShortName}
type ApiService struct {
}

//GetServices 获取当前插件提供的所有服务
func (p *ApiService) GetServices() []string {
	return GetServices()
}

//Handle 业务处理
func (p *ApiService) Handle(name string, mode string, service string, ctx *context.Context) (response context.Response, err error) {
	response, err = Handle(name, mode, service, ctx)
	if err != nil {
		err = fmt.Errorf("{@pImportName}:status:%d,err:%v", response.GetStatus(err), err)
	}
	return
}
func (p *ApiService) Close() error {
	return nil
}

//GetWorker 获取当前worker
func GetWorker() context.Worker {
	return &ApiService{}
}
`
var TmplMap map[string]string

func init() {
	TmplMap = map[string]string{
		"main.go":               mainTmpl,
		"registry.go":           routerTmpl,
		"services/myservice.go": handlers.HandlerOrderTmpl,
		"conf/conf.go":          conf.ConfTmpl,
	}
}
