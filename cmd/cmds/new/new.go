package new

import (
	"errors"
	"os"

	"fmt"

	"path/filepath"

	"strings"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/gaea/cmd/cmds/hydra"
	"github.com/qxnw/gaea/cmd/cmds/new/api"
	"github.com/qxnw/gaea/cmd/cmds/new/web"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/transform"
	"github.com/spf13/pflag"
)

//command 根据输入参数打包项目
type command struct {
	logger      *logger.Logger
	projectName string
	newWeb      bool
	tf          *transform.Transform
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) error {
	flags.BoolVarP(&r.newWeb, "web project", "w", false, "指定为web项目")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if len(os.Args) < 2 {
		return errors.New("未指定目路径")
	}
	r.projectName = os.Args[2]
	if r.projectName == "" {
		return errors.New("未指定目路径")
	}
	r.tf = r.makeParams()
	return nil
}

//Run 运行指令
func (r *command) Run(args []string) error {
	fullName, err := hydra.GetProjectSrcPath(r.projectName)
	if err != nil {
		return err
	}
	if _, err := os.Stat(fullName); !os.IsNotExist(err) {
		err = errors.New("项目已存在或不为空")
		return err
	}
	if r.newWeb {
		err = r.createProject(fullName, web.TmplMap)
	} else {
		err = r.createProject(fullName, api.TmplMap)
	}
	if err != nil {
		return err
	}
	r.logger.Info("项目生成成功:", fullName)
	return nil

}
func (r *command) createProject(root string, data map[string]string) error {
	for k, v := range data {
		path := filepath.Join(root, k)
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			err = fmt.Errorf("创建文件夹%s失败:%v", path, err)
			return err
		}
		srcf, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
			return err
		}
		defer srcf.Close()
		n, err := srcf.WriteString(r.tf.Translate(v))
		if err != nil {
			return err
		}
		r.logger.Info("创建文件:", path, n)
	}
	return nil

}
func (r *command) makeParams() *transform.Transform {
	names := strings.Split(strings.Trim(r.projectName, "/"), "/")
	className := ""
	for _, v := range names {
		className += strings.ToUpper(v[0:1]) + v[1:]
	}
	tf := transform.New()
	tf.Set("pShortName", names[len(names)-1])
	tf.Set("pImportName", strings.Trim(r.projectName, "/"))
	tf.Set("pClassName", className)
	return tf
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("new", &commandResolver{})
}
