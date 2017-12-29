package build

import (
	"os"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/gaea/cmd/cmds/hydra"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

//command 根据输入参数打包项目
type command struct {
	logger   *logger.Logger
	projects []string
	fileName string
	version  string
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) error {
	flags.StringVarP(&r.version, "项目版本号", "v", "", "项目版本号")
	flags.Parse(os.Args[1:])
	r.projects = os.Args[2:]
	if len(r.version) > 0 {
		r.projects = os.Args[4:]
	} else {
		r.version = "0.0.1"
	}

	return nil
}

//Run 运行指令
func (r *command) Run(args []string) (err error) {
	err = r.BuildProjects(r.projects)
	if err != nil {
		return
	}
	r.logger.Info("开始编译:", "hydra")	
	err = hydra.BuildHydra(r.version)
	if err != nil {
		return
	}
	r.logger.Info("编译成功")
	return nil
}
func (r *command) BuildProjects(projectNames []string) error {
	for _, p := range projectNames {
		r.logger.Info("开始编译:", p)
		err := hydra.BuildPlugin(p)
		if err != nil {
			return err
		}
	}
	return nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("build", &commandResolver{})
}
