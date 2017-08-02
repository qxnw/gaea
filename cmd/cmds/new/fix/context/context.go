package context

var ContextTmpl = `package context

import (

	"errors"
	"sync"

	"github.com/qxnw/goplugin"
	"github.com/qxnw/goplugin/errorCode"
	"github.com/qxnw/lib4go/encoding/base64"
	"github.com/qxnw/lib4go/security/xsrf"
	"github.com/qxnw/lib4go/utility"
	
	"{@pImportName}/conf"

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

func GetContext(ctx goplugin.Context, rpc goplugin.RPCInvoker, fields map[string][]string) (cx *Context, status int, p map[string]interface{}, err error) {
	context, err := goplugin.GetContext(ctx, rpc)
	if err != nil {
		status = errorCode.SERVER_ERROR
		return
	}

	err = context.CheckInput(fields["input"]...)
	if err != nil {
		status = errorCode.NOT_ACCEPTABLE
		context.Close()
		return
	}
	err = context.CheckArgs(fields["args"]...)
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

//Redirect 页面转跳
func (c *Context) Redirect(code int, url string) {
	c.PluginContext.Redirect(code, url, c.Param)
}

//SetView 设置view
func (c *Context) SetView(name string) {
	c.Param["__view"] = name
}

//NoView 设置view
func (c *Context) NoView() {
	c.Param["__view"] = "NONE"
}

//GetXSRFToken
func (c *Context) GetXSRFToken() string {
	t, _ := base64.URLDecode(utility.GetGUID())
	token := xsrf.CreateXSRFToken(c.Conf.XSRFSecret, t)
	return token
}

//GetSessionID 获取当session
func (c Context) GetUserSessionID() string {
	return c.GetCookieString("{@pShortName}_sid")
}

//SetSessionID 设置session
func (c *Context) SetSessionID(sessionID string) string {
	c.SetCookie("{@pShortName}_sid", sessionID, 3600*24*15, c.Conf.HostName, c.Param)
	return sessionID
}

func (c *Context) IsLogin() (err error) {
	sid := c.GetUserSessionID()
	if sid == "" {
		err = errors.New("sid为空，判定为session过期或新用户，需要跳转授权")
		return
	}
	return
}


func (c *Context) Close() {
	contextPool.Put(c)
	if c.PluginContext != nil {
		c.PluginContext.Close()
	}
}
`
