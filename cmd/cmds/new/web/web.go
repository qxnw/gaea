package web

import (
	"github.com/qxnw/gaea/cmd/cmds/new/web/conf"
	"github.com/qxnw/gaea/cmd/cmds/new/web/context"
	"github.com/qxnw/gaea/cmd/cmds/new/web/controllers"
	"github.com/qxnw/gaea/cmd/cmds/new/web/views"
)

var mainTmpl = `package main

//go build -buildmode=plugin
import (
	"fmt"

	"github.com/qxnw/goplugin"
)

type {@pShortName} struct {
}

//GetServices 获取当前插件提供的所有服务
func (p *{@pShortName}) GetServices() []string {
	return GetServices()
}

//Handle 业务处理
func (p *{@pShortName}) Handle(name string, mode string, service string, c goplugin.Context, rpc goplugin.RPCInvoker) (status int, result interface{}, param map[string]interface{}, err error) {
	if h, ok := GetHandlers()[service]; ok {
		status, r, param, err := h.Handle(service, c, rpc)
		if err != nil || status != 200 {
			return status, result, nil, fmt.Errorf("{@pImportName}:status:%d,err:%v", status, err)
		}
		return status, r, param, err
	}
	return 404, "", nil, fmt.Errorf("{@pImportName} 未找到服务:%s", service)
}
func (p *{@pShortName}) Close() error {
	return nil
}

//GetWorker 获取当前worker
func GetWorker() goplugin.Worker {
	return &{@pShortName}{}
}`

var TmplMap map[string]string

//{@pShortName} xwuser
//{@pClassName} XwUser
//{@pImportName} /xwll/xwuser

func init() {
	TmplMap = map[string]string{
		"main.go":              mainTmpl,
		"routers.go":            routerTmpl,
		"views/index.html":     views.IndexViewTmpl,
		"controllers/index.go": controllers.IndexControllerTmpl,
		"context/context.go":   context.ContextTmpl,
		"conf/conf.go":         conf.ConfTmpl,
	}
}
