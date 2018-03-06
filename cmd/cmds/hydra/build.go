package hydra

import (
	"fmt"
	"os"
	"os/exec"

	"strings"
	"sync"
)

var buildLock sync.Mutex

func BuildHydra(version string) error {
	return goInstallHydra("github.com/qxnw/hydra", version)
}

// goInstall 调用go install生成项目程序
func goInstallHydra(projectShortName string, version string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := GetHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)
	os.Remove("hydra")
	if version != "" {
		icmd := exec.Command("go", "install", "-ldflags", fmt.Sprintf("-X main.VERSION=%s", version), strings.Trim(projectShortName, "//"))
		icmd.Stdout = os.Stdout
		icmd.Stderr = os.Stderr
		return icmd.Run()
	}
	icmd := exec.Command("go", "install", strings.Trim(projectShortName, "//"))
	icmd.Stdout = os.Stdout
	icmd.Stderr = os.Stderr
	return icmd.Run()
}

func GoBuild(projectShortName string, version string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := GetHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)
	if version != "" {
		icmd := exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.VERSION=%s", version), strings.Trim(projectShortName, "//"))
		_, err = icmd.Output()
		fmt.Println("err.x:", err, strings.Trim(projectShortName, "//"))
		if err != nil && strings.Contains(err.Error(), "target main.main not defined") {
			return nil
		}
		return err
	}
	icmd := exec.Command("go", "build", strings.Trim(projectShortName, "//"))
	_, err = icmd.Output()
	fmt.Println("err.x:", err, strings.Trim(projectShortName, "//"))
	if err != nil && strings.Contains(err.Error(), "target main.main not defined") {
		return nil
	}
	return err
}

//BuildPlugin  go build -buildmode=plugin
func BuildPlugin(projectShortName string, version string) error {
	buildLock.Lock()
	defer buildLock.Unlock()

	workDir, err := GetHydraExecDir()
	if err != nil {
		return err
	}
	os.Chdir(workDir)
	if version != "" {
		icmd := exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.VERSION=%s", version), "-buildmode=plugin", strings.Trim(projectShortName, "//"))
		icmd.Stdout = os.Stdout
		icmd.Stderr = os.Stderr
		icmd.Env = append(os.Environ(), "GOGC=off")
		return icmd.Run()
	}
	icmd := exec.Command("go", "build", "-buildmode=plugin", strings.Trim(projectShortName, "//"))
	icmd.Stdout = os.Stdout
	icmd.Stderr = os.Stderr
	icmd.Env = append(os.Environ(), "GOGC=off")
	return icmd.Run()
}
