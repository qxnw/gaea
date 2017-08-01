package context

var ContextTmpl = `package context
import (
	"github.com/qxnw/goplugin"
	"github.com/qxnw/goplugin/errorCode"
	"{@pImportName}/conf"
	"sync"
)
var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}
}

type Context struct {
	*goplugin.PluginContext
	Param map[string]interface{}
	Conf  *conf.{@pClassName}Conf
}
func GetContext(ctx goplugin.Context, invoker goplugin.RPCInvoker, service map[string][]string) (cx *Context, status int, p map[string]interface{}, err error) {

	context, err := goplugin.GetContext(ctx, invoker)
	if err != nil {
		status = errorCode.SERVER_ERROR
		return
	}
	err = context.CheckInput(service["mustFields"]...)
	if err != nil {
		status = errorCode.NOT_ACCEPTABLE
		context.Close()
		return
	}
	err = context.CheckArgs(service["mustArgs"]...)
	if err != nil {
		status = errorCode.NOT_EXTENDED
		context.Close()
		return
	}
cn, err := conf.GetConf(context)
	if err != nil {
		status = errorCode.SERVER_ERROR
		context.Close()
		return
	}

	cx = contextPool.Get().(*Context)
	cx.PluginContext = context
	cx.Param = make(map[string]interface{})
	cx.Conf = cn
	p = cx.Param
	return
}

func (c *Context) Close() {
	contextPool.Put(c)
	if c.PluginContext != nil {
		c.PluginContext.Close()
	}
}
`
