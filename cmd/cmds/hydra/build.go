package hydra

import (
	"fmt"
	"os"
	"os/exec"

	"strings"
	"sync"
)

var buildLock sync.Mutex

func buildHydra() error {
	return goInstall("github.com/qxnw/hydra")
}

// goInstall 调用go install生成项目程序
func goInstall(projectShortName string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := getHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)

	icmd := exec.Command("go", "install", strings.Trim(projectShortName, "//"))
	icmd.Stdout = os.Stdout
	icmd.Stderr = os.Stderr
	icmd.Env = append(os.Environ(), "GOGC=off")
	return icmd.Run()
}

func goBuild(projectShortName string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := getHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)
	icmd := exec.Command("go", "build", strings.Trim(projectShortName, "//"))
	_, err = icmd.Output()
	fmt.Println("err.x:", err, strings.Trim(projectShortName, "//"))
	if err != nil && strings.Contains(err.Error(), "target main.main not defined") {
		return nil
	}
	return err
}

//goBuild  go build -buildmode=plugin
func goBuildPlugin(projectShortName string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := getHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)

	icmd := exec.Command("go", "build", "-buildmode=plugin", strings.Trim(projectShortName, "//"))
	icmd.Stdout = os.Stdout
	icmd.Stderr = os.Stderr
	icmd.Env = append(os.Environ(), "GOGC=off")
	return icmd.Run()
}
