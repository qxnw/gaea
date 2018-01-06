package pack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/security/crc32"
	"github.com/qxnw/lib4go/utility"
	"github.com/spf13/pflag"
	"github.com/zkfy/archiver"
)

//command 根据输入参数打包项目
type command struct {
	logger    *logger.Logger
	packFiles []string
	fileName  string
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) error {
	flags.StringVarP(&r.fileName, "打包文件名称", "n", "", "打包文件名称")

	flags.Parse(os.Args[1:])
	if len(os.Args) < 3 {
		return errors.New("未指定打包文件")
	}
	if len(r.fileName) > 0 {
		r.packFiles = os.Args[4:]
	} else {
		r.packFiles = os.Args[2:]
	}

	if len(r.fileName) == 0 {
		r.fileName = utility.GetGUID()[:6]
	}
	if !strings.HasSuffix(r.fileName, ".tar.gz") {
		r.fileName = fmt.Sprintf("./%s.tar.gz", r.fileName)
	}
	return nil
}

//Run 运行指令
func (r *command) Run(args []string) (err error) {

	err = r.PackProjects(r.fileName, r.packFiles)
	if err != nil {
		r.logger.Errorf("打包失败:%v", err)
		return err
	}
	buff, err := ioutil.ReadFile(r.fileName)
	if err != nil {
		err = fmt.Errorf("读取文件失败:%v", err)
		return
	}

	r.logger.Info("打包完成:", r.fileName, crc32.Encrypt(buff))

	return nil
}
func (r *command) PackProjects(dest string, projectName []string) error {
	packFiles := []string{"hydra"}
	for _, p := range projectName {
		names := strings.Split(p, "/")
		packFiles = append(packFiles, names[len(names)-1]+".so")
		p, err := getProjectPath(p)
		if err != nil {
			return err
		}

		path := filepath.Join(p, "views")
		if _, err := os.Stat(path); err == nil {
			packFiles = append(packFiles, path)
		}
		static := filepath.Join(p, "static")
		if _, err := os.Stat(static); err == nil {
			packFiles = append(packFiles, static)
		}
	}

	return r.PackFiles(dest, packFiles)
}
func (r *command) PackFiles(dest string, files []string) error {
	for _, f := range files {
		if _, ex := os.Stat(f); os.IsNotExist(ex) {
			err := fmt.Errorf("文件不存在:%s err:%v", f, ex)
			return err
		}
	}
	r.logger.Info("打包文件:", files)
	err := archiver.TarGz.Make(dest, files)
	if err != nil {
		if _, ex := os.Stat(dest); ex == nil {
			os.Remove(dest)
		}
		return err
	}
	return nil
}
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
func getProjectPath(pName string) (path string, err error) {
	gopath, err := getGoPath()
	if err != nil {
		return "", err
	}
	for _, v := range gopath {
		srcPath := filepath.Join(v, "/src/")
		path := filepath.Join(srcPath, pName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("pack", &commandResolver{})
}
