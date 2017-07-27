package cmds

import (
	"fmt"

	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

//ICommand 通用命令执行工具
type ICommand interface {
	PreRun(*pflag.FlagSet) error
	Run([]string) error
}

//IWatcherResolver 定义命令工具生成方法
type IWatcherResolver interface {
	Resolve(name string, log *logger.Logger) (ICommand, error)
}

var commands map[string]IWatcherResolver

func init() {
	commands = make(map[string]IWatcherResolver)
}

//Register 注册命令工具
func Register(name string, cmd IWatcherResolver) {
	if _, ok := commands[name]; ok {
		panic("已经包含命令：" + name)
	}
	commands[name] = cmd
}

//NewCommand 返回命令执行工具
func NewCommand(name string, log *logger.Logger) (ICommand, error) {
	cmds, ok := commands[name]
	if !ok {
		return nil, fmt.Errorf("未找到命令: %q (forgotten import?)", name)
	}
	return cmds.Resolve(name, log)
}
