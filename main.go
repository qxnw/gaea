package main

import (
	"os"

	"github.com/qxnw/gaea/cmd"
	_ "github.com/qxnw/gaea/cmd/cmds/hydra"
	_ "github.com/qxnw/gaea/cmd/cmds/new"
	_ "github.com/qxnw/gaea/cmd/cmds/pack"
	"github.com/qxnw/lib4go/logger"
	"github.com/spf13/pflag"
)

func main() {
	logger := logger.GetSession("gaea", logger.CreateSession())
	logger.Info("启动 gaea ...")
	defer func() {
		logger.WaitClose()
	}()
	if len(os.Args) < 2 {
		logger.Error("未指定命令名称：run,pack,hydra")
		return
	}
	name := os.Args[1]
	cmd, err := cmds.NewCommand(name, logger)
	if err != nil {
		logger.Error(err)
		return
	}
	err = cmd.PreRun(pflag.CommandLine)
	if err != nil {
		logger.Error(err)
		return
	}
	err = cmd.Run(pflag.Args())
	if err != nil {
		logger.Error(err)
		return
	}

}
