package pack

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"path"

	"path/filepath"

	"io"

	"strings"

	"fmt"

	"github.com/qxnw/gaea/cmd"
	"github.com/qxnw/gaea/cmd/cmds/hydra"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/utility"
	"github.com/spf13/pflag"
)

//command 根据输入参数打包项目
type command struct {
	logger      *logger.Logger
	projectName string
	packViews   bool
}

//PreRun 预执行用于绑定输入参数及运行前初始化
func (r *command) PreRun(flags *pflag.FlagSet) error {
	flags.StringVarP(&r.projectName, "打包项目", "p", "", "需要打包的域或项目名称")
	flags.BoolVarP(&r.packViews, "打包views", "v", false, "打包web项目的views文件")
	flags.Parse(os.Args[1:])
	if r.projectName == "" {
		flags.Usage()
		return errors.New("未指定打包项目路径")
	}
	return nil
}

//Run 运行指令
func (r *command) Run(args []string) error {
	subNames, _, err := hydra.GetSubSystemSrcDir(r.projectName)
	if err != nil {
		return err
	}
	root, pkg, bin, views, err := r.createPackPath()
	if err != nil {
		return err
	}
	r.logger.Info("编译hydra")
	err = hydra.BuildHydra()
	if err != nil {
		err = errors.New("hydra编译失败")
		return err
	}
	err = r.copy(root, bin, "hydra")
	if err != nil {
		return err
	}
	for _, name := range subNames {
		r.logger.Infof("编译插件:%s", name)
		err = hydra.BuildPlugin(name)
		if err != nil {
			r.logger.Errorf("插件:%s编译失败", name)
			return err
		}
		r.logger.Errorf("插件:%s编译成功", name)
		err = r.copy(root, bin, r.getPluginName(name))
		if err != nil {
			r.logger.Errorf("复制插件:%s文件失败", name)
			return err
		}
	}
	if r.packViews {
		for _, name := range subNames {
			pname, err := hydra.GetProjectPath(name)
			if err != nil {
				continue
			}
			pn := filepath.Join(pname, "views")
			if _, err = os.Stat(pn); err == nil {
				filepath.Walk(pn, func(filename string, fi os.FileInfo, err error) error {
					r.copyFile(filename, filepath.Join(views, path.Base(filename)))
					return nil
				})
			}
		}
	}

	return nil
}
func (r *command) getPluginName(n string) string {
	names := strings.Split(n, "/")
	return names[len(names)-1]
}
func (r *command) createPackPath() (root string, packRoot string, bin string, views string, err error) {
	root, err = hydra.GetHydraExecDir()
	if err != nil {
		return "", "", "", "", err
	}
	packRoot = filepath.Join(root, utility.GetGUID())
	err = os.Mkdir(packRoot, 777)
	if err != nil {
		return "", "", "", "", err
	}
	bin = filepath.Join(packRoot, "/bin")
	err = os.Mkdir(bin, 777)
	if err != nil {
		return "", "", "", "", err
	}
	if !r.packViews {
		return "", "", "", "", nil
	}
	views = filepath.Join(packRoot, "/views")
	err = os.Mkdir(views, 777)
	if err != nil {
		return "", "", "", "", err
	}
	return "", "", "", "", nil
}
func (r *command) copy(src, dst string, name string) error {
	src = filepath.Join(src, name)
	dst = filepath.Join(dst, name)
	return r.copyFile(src, dst)
}
func (r *command) copyFile(src, dst string) error {
	srcf, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("无法打开文件:%s(err:%v)", src, err)
		return err
	}
	defer srcf.Close()
	dstDir := filepath.Dir(dst)
	err = os.MkdirAll(dstDir, 777)
	if err != nil {
		err = fmt.Errorf("创建文件夹失败:%s(err:%v)", dstDir, err)
		return err
	}
	dstf, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("创建文件失败:%s(err:%v)", dst, err)
		return err
	}
	defer dstf.Close()
	_, err = io.Copy(dstf, srcf)
	return err
}
func (r *command) createZIP(files []string, zipName string) error {
	buf := new(bytes.Buffer)

	// 创建一个压缩文档
	w := zip.NewWriter(buf)
	for _, file := range files {
		f, err := w.Create(path.Base(file))
		if err != nil {
			return err
		}
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fi.Close()
		_, err = io.Copy(f, fi)
		if err != nil {
			return err
		}
	}
	err := w.Close()
	if err != nil {
		return err
	}
	// 将压缩文档内容写入文件
	f, err := os.OpenFile(zipName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}

type commandResolver struct {
}

func (r *commandResolver) Resolve(name string, log *logger.Logger) (cmds.ICommand, error) {
	return &command{logger: log}, nil
}
func init() {
	cmds.Register("pack", &commandResolver{})
}
