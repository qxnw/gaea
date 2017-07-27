package hydra

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/qxnw/lib4go/logger"
)

type hydraProcess struct {
	process map[string]*os.Process
	logger  *logger.Logger
}

func (p *hydraProcess) Restart(params []string, projectName ...string) {
	err := p.preBuild()
	if err != nil {
		return
	}
	var buildSync sync.WaitGroup
	buildSync.Add(2)
	var errBuild error
	var errKill error
	go func() {
		for _, v := range projectName {
			errBuild = p.buildPlugin(v)
			if errBuild != nil {
				break
			}
		}
		buildSync.Done()
	}()
	go func() {
		errKill = p.kill(params)
		time.Sleep(time.Second)
		buildSync.Done()
	}()

	buildSync.Wait()
	if errBuild == nil && errKill == nil {
		go p.startHydra(params)
	}
}

func (w *hydraProcess) preBuild() error {
	w.logger.Info("开始编译hydra...")
	err := goInstall("github.com/qxnw/hydra")
	if err != nil {
		w.logger.Error("hydra编译失败:", err)
		return err
	}
	w.logger.Info("hydra编译成功")
	return nil
}
func (w *hydraProcess) kill(runParam []string) error {
	w.logger.Info("结束进程:hydra ", runParam)
	err := kill(runParam)
	if err != nil {
		w.logger.Error("进程未结束：hydra", err)
	}
	w.logger.Info("hydra已关闭")
	return err
}
func (w *hydraProcess) buildPlugin(projectName string) error {
	w.logger.Infof("开始编译%s...", projectName)
	err := goBuildPlugin(projectName)
	if err != nil {
		w.logger.Error(projectName, "编译失败:", err)
		return err
	}
	w.logger.Info(projectName, "编译成功")
	return nil
}
func (w *hydraProcess) startHydra(runParam []string) error {
	w.logger.Infof("-------------开始启动hydra-------------")
	err := start(runParam)
	if err != nil {
		w.logger.Error("启动失败:", err)
		return err
	}
	return nil
}

func (p *hydraProcess) Start(params []string) error {
	workDir, err := getHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)

	icmd := exec.Command("./hydra", params...)
	icmd.Stdin = os.Stdin
	icmd.Stdout = os.Stdout
	runChan := make(chan error, 1)
	go func() {
		err := icmd.Run()
		runChan <- err
	}()

	for {
		select {
		case err = <-runChan:
			return err
		case <-time.After(time.Second * 3):
			key := strings.Join(params, "|")
			p.process[key] = icmd.Process
			return nil
		}
	}

}
func (p *hydraProcess) Kill(params []string) error {
	key := strings.Join(params, "|")
	if pc, ok := p.process[key]; ok {
		return pc.Signal(os.Interrupt)
	}
	return nil
}

var process = &hydraProcess{process: make(map[string]*os.Process), logger: logger.GetSession("process", logger.CreateSession())}

func start(params []string) error {
	return process.Start(params)
}
func kill(params []string) error {
	return process.Kill(params)
}

func restart(params []string, projectName ...string) {
	process.Restart(params, projectName...)
}
