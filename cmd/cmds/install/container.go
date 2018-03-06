package install

import (
	"fmt"

	"github.com/qxnw/hydra/registry"
	"github.com/qxnw/lib4go/logger"
)

func GetVarHandler(domain string, addr string, log logger.ILogger) func(tp string, name string) (string, error) {
	return func(tp string, name string) (string, error) {
		registry, err := registry.NewRegistryWithAddress(addr, log)
		if err != nil {
			return "", err
		}
		buff, _, err := registry.GetValue(fmt.Sprintf("/%s/var/%s/%s", domain, tp, name))
		if err != nil {
			return "", err
		}
		return string(buff), nil
	}
}
