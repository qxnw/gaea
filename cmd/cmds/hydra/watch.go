package hydra

import (
	"time"

	"os"

	"strings"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/qxnw/lib4go/logger"
)

var (
	ignoreFile = []string{
		`.#(\w+).go`,
		`.(\w+).go.swp`,
		`(\w+).go~`,
		`(\w+).tmp`,
	}
	watchFiles = []string{".go"}
	eventTime  = make(map[string]int64)
)

//Watcher 文件监控器
type Watcher struct {
	watcher     *fsnotify.Watcher
	projectName string
	logger      *logger.Logger
	closeCh     chan struct{}
	lastBuild   time.Time
	runParam    []string
	buildSync   sync.Mutex
	isBuilding  bool
}

//NewWatcher 创建项目文件监控
func NewWatcher(projectName string, runParam []string, paths []string, logger *logger.Logger) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		err = watcher.Add(path)
		if err != nil {
			return nil, err
		}
	}
	return &Watcher{watcher: watcher,
		projectName: projectName,
		logger:      logger,
		runParam:    runParam,
		lastBuild:   time.Now()}, nil
}

//Start 启动文件监控
func (w *Watcher) Start() {
	go func() {
	LOOP:
		for {
			select {
			case <-w.closeCh:
				break LOOP
			case e := <-w.watcher.Events:
				isBuild := true
				if !isCurentFile(e.Name, watchFiles) || isCurentFile(e.Name, ignoreFile) {
					continue
				}
				mt := getFileModTime(e.Name)
				if t := eventTime[e.Name]; mt == t {
					isBuild = false
				}
				eventTime[e.Name] = mt

				if isBuild && time.Since(w.lastBuild).Seconds() > 2 && !w.isBuilding {
					w.buildSync.Lock()
					w.isBuilding = !w.isBuilding
					w.buildSync.Unlock()
					if w.isBuilding {
						w.lastBuild = time.Now()
						go func() {
							w.logger.Info(w.projectName, "文件发生变化，启动重新编译...")
							time.Sleep(time.Second * 1)
							restart(w.runParam, w.projectName)
							w.buildSync.Lock()
							w.isBuilding = false
							w.buildSync.Unlock()
						}()
					}

				}
			case err := <-w.watcher.Errors:
				w.logger.Warnf("Watcher error: %s", err.Error()) // No need to exit here
			}
		}
		w.watcher.Close()
	}()
}

//Close 关闭监控
func (w *Watcher) Close() {
	close(w.closeCh)
}
func getFileModTime(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Now().Unix()
	}
	return fi.ModTime().Unix()
}
func isCurentFile(path string, exts []string) bool {
	for _, v := range exts {
		if strings.HasSuffix(path, v) {
			return true
		}
	}
	return false
}
