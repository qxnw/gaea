package pack

import (
	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

//command 根据输入参数打包项目
type command struct {
	logger *logger.Logger
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) error {
	return nil
}

//Run 运行指令
func (r *command) Run(args []string) error {
	return nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("pack", &commandResolver{})
}
