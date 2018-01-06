package install

import (
	"errors"
	"fmt"
	"os"
	"plugin"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/hydra/context"
	"github.com/qxnw/lib4go/file"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

//command 用于安装插件中指定的数据库，注册中心等
type command struct {
	logger          *logger.Logger
	registryAddress string
	domain          string
	name            string
}

//PreRun 预执行用于绑定输入参数及运行前初始化 gaea install /hydra -r "zk://192.168.0.159" hydra_core.so
func (r *command) PreRun(flags *pflag.FlagSet) error {
	flags.StringVarP(&r.registryAddress, "registry center address", "r", "", "注册中心地址(格式：zk://192.168.0.159:2181,192.168.0.158:2181)")
	if len(os.Args) < 6 {
		flags.Usage()
		return errors.New("未指定域或注册中心")
	}
	flags.Parse(os.Args[1:])
	r.domain = os.Args[2]
	r.name = os.Args[5]
	return nil
}

//Run 运行指令
func (r *command) Run(args []string) (err error) {
	r.logger.Info("加载:", r.name)
	f, err := r.loadPlugin(r.name)
	if err != nil {
		return
	}
	r.logger.Info("开始安装")

	err = f(r.domain, GetVarHandler(r.domain, r.registryAddress, r.logger))
	if err != nil {
		return
	}
	r.logger.Info("安装成功")
	return err
}
func (r *command) loadPlugin(p string) (f func(domain string, r context.VarHandle) error, err error) {
	path, err := file.GetAbs(p)
	if err != nil {
		return
	}
	if _, err = os.Lstat(path); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("%s不存在", p)
	}

	pg, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("加载插件失败:%s,err:%v", path, err)
	}
	work, err := pg.Lookup("Install")
	if err != nil {
		return nil, fmt.Errorf("加载%s失败未找到函数Install,err:%v", path, err)
	}
	wkr, ok := work.(func(domain string, h context.VarHandle) error)
	if !ok {
		return nil, fmt.Errorf("加载%s失败 Install函数必须为func(domain string, r context.VarHandle) error类型", path)
	}
	return wkr, nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("install", &commandResolver{})
}
