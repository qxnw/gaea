package api

var mainTmpl = `
package main

import (
	"@pImportName/services/order"
	"github.com/qxnw/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("api"),
		hydra.WithAutoCreateConf(),
		hydra.WithDebug())

	app.Micro("/order/query", order.NewQueryHandler)

	app.Start()
}

`
var TmplMap map[string]string

func init() {
	TmplMap = map[string]string{
		"main.go":                       mainTmpl,
		"services/order/order.query.go": orderQueryTmpl,
	}
}
