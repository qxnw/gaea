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
		return nil, fmt.Errorf("环境变量未配置gopath")
	}
	path := strings.Split(gopath, ";")
	if len(path) == 0 {
		return nil, fmt.Errorf("环境变量gopath配置的路径为空")
	}
	return path, nil
}

func getHydraSrcDir() (string, error) {
	return getProjectPath("github.com/qxnw/hydra")
}
func getProjectPath(short string) (string, error) {
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
func getHydraExecDir() (string, error) {
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
	return "", fmt.Errorf("未找到hydra执行文件:%v", gopath)
}
func getSubSystemSrcDir(groupName string) ([]string, []string, error) {
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
