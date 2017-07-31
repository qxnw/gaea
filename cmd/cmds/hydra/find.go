package hydra

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getGoPath() ([]string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = os.Getenv("HOME")
		if gopath != "" {
			gopath = filepath.Join(gopath, "work")
		}
	}
	if gopath == "" {
		return nil, fmt.Errorf("未配置环境变量GOPATH")
	}
	path := strings.Split(gopath, ";")
	if len(path) == 0 {
		return nil, fmt.Errorf("环境变量GOPATH配置的路径为空")
	}
	return path, nil
}

func GetHydraSrcDir() (string, error) {
	return GetProjectPath("github.com/qxnw/hydra")
}
func GetProjectPath(short string) (string, error) {
	gopath, err := getGoPath()
	if err != nil {
		return "", err
	}
	for _, v := range gopath {
		path := filepath.Join(v, "/src/", short)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到项目文件:%v", short)
}
func GetHydraExecDir() (string, error) {
	gopath, err := getGoPath()
	if err != nil {
		return "", err
	}
	for _, v := range gopath {
		path := filepath.Join(v, "/bin/hydra")
		if _, err := os.Stat(path); err == nil {
			return filepath.Join(v, "/bin"), nil
		}
	}
	return filepath.Join(gopath[0], "/bin"), nil
}
func GetSubSystemSrcDir(groupName string) ([]string, []string, error) {
	gopath, err := getGoPath()
	if err != nil {
		return nil, nil, err
	}
	proShortNames := make([]string, 0, 2)
	profullNames := make([]string, 0, 2)
	for _, v := range gopath {
		srcPath := filepath.Join(v, "/src/")
		path := filepath.Join(srcPath, groupName)
		if _, err := os.Stat(path); err == nil {
			filepath.Walk(path, func(filename string, fi os.FileInfo, err error) error {
				if strings.HasSuffix(filename, "main.go") && isPluginFile(filename) {
					dir := filepath.Dir(filename)
					proShortNames = append(proShortNames, dir[len(srcPath):])
					profullNames = append(profullNames, dir)
				}
				return nil
			})
		}
	}
	return proShortNames, profullNames, nil
}
func isPluginFile(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			return false
		}
		if strings.Contains(line, "GetWorker") && strings.Contains(line, "goplugin.Worker") {
			return true
		}
	}
}
