package hydra

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"path/filepath"

	"strings"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/hydra/hydra"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

//command 根据输入参数打包项目
type command struct {
	logger            *logger.Logger
	subSysSrcDirNames []string
	subSysSrcDirs     []string
	watchers          []*Watcher
	dependencePath    []string
	*hydra.HFlags
	domainServer string
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) (err error) {
	r.HFlags.BindFlags(flags)
	flags.StringVarP(&r.domainServer, "监控的服务器", "w", "", "监控服务器域或服务器名称")

	err = r.HFlags.CheckFlags(2)
	if err != nil {
		return err
	}
	if r.domainServer == "" {
		r.domainServer = r.Domain
	}

	if r.domainServer == "" {
		pflag.Usage()
		return fmt.Errorf("缺少输入参数:domain")
	}
	r.watchers = make([]*Watcher, 0, 1)
	r.subSysSrcDirNames, r.subSysSrcDirs, err = GetSubSystemSrcDir(r.domainServer)
	if err != nil {
		return
	}
	if len(r.subSysSrcDirNames) == 0 {
		path, err := getGoPath()
		if err != nil {
			return err
		}
		err = fmt.Errorf("在目录：%v的src目录下未找到:%s的项目文件", r.domainServer, path)
		return err
	}
	hydra, err := GetProjectPath("github.com/qxnw/hydra")
	if err != nil {
		return err
	}
	lib4go, err := GetProjectPath("github.com/qxnw/lib4go")
	if err != nil {
		return err
	}
	goPlugin, err := GetProjectPath("github.com/qxnw/goplugin")
	if err != nil {
		return err
	}
	r.dependencePath = append(r.dependencePath, r.getChildrenDirs(hydra)...)
	r.dependencePath = append(r.dependencePath, r.getChildrenDirs(lib4go)...)
	r.dependencePath = append(r.dependencePath, r.getChildrenDirs(goPlugin)...)

	return nil
}

//Run 运行指令
func (r *command) Run(args []string) error {
	for i, v := range r.subSysSrcDirs {
		paths := r.getChildrenDirs(v)
		paths = append(paths, r.dependencePath...)
		r.logger.Infof("监控项目:%s文件的变化...%d-%d-%v", v, len(paths), len(args), args)
		watcher, err := NewWatcher(r.subSysSrcDirNames[i], r.HFlags.ToArgs(), paths, r.logger)
		if err != nil {
			return err
		}
		watcher.Start()
		r.watchers = append(r.watchers, watcher)
	}
	//启动项目
	restart(r.HFlags.ToArgs(), r.subSysSrcDirNames...)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM) //9:kill/SIGKILL,15:SIGTEM,20,SIGTOP 2:interrupt/syscall.SIGINT
LOOP:
	for {
		select {
		case <-interrupt:
			r.logger.Warnf("gaea (%s) was killed", r.domainServer)
			kill(r.HFlags.ToArgs())
			break LOOP
		}
	}
	return nil
}
func (r *command) getChildrenDirs(dir string) []string {
	paths := make([]string, 0, 8)
	filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		if fi != nil && fi.IsDir() && !strings.Contains(filename, ".git") {
			paths = append(paths, filename)
		}
		return nil
	})
	return paths
}
func (r *command) Close() error {
	for _, v := range r.watchers {
		v.Close()
	}
	return nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log, HFlags: &hydra.HFlags{}}, nil
}
func init() {
	cmds.Register("hydra", &commandResolver{})
}
