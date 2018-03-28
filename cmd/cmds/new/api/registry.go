package api

var registryTmpl = `
package main

import (
	"github.com/qxnw/hydra/component"
	"github.com/qxnw/hydra/engines"
	"@pImportName/services/order"
)

func loader() engines.ServiceLoader {
	return func(component *component.StandardComponent, container component.IContainer) error {	
		component.AddMicroService("/order/query", order.NewQueryHandler)
		return nil
	}
}


`
