package main

import (
	"os"

	"github.com/qxnw/gaea/cmd"
	_ "github.com/qxnw/gaea/cmd/cmds/build"
	_ "github.com/qxnw/gaea/cmd/cmds/hydra"
	_ "github.com/qxnw/gaea/cmd/cmds/install"
	_ "github.com/qxnw/gaea/cmd/cmds/new"
	_ "github.com/qxnw/gaea/cmd/cmds/pack"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

func main() {
	gaeaLogger := logger.GetSession("gaea", logger.CreateSession())

	defer logger.Close()

	if len(os.Args) < 2 {
		gaeaLogger.Error("未指定命令名称：new,run,pack,build,install")
		return
	}
	name := os.Args[1]
	cmd, err := cmds.NewCommand(name, gaeaLogger)
	if err != nil {
		gaeaLogger.Error(err)
		return
	}
	err = cmd.PreRun(pflag.CommandLine)
	if err != nil {
		gaeaLogger.Error(err)
		return
	}
	gaeaLogger.Info("启动 gaea ...")
	err = cmd.Run(pflag.Args())
	if err != nil {
		gaeaLogger.Error(err)
		return
	}

}
