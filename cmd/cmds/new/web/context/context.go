package context

var ContextTmpl = `package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"{@pImportName}/conf"

	"github.com/qxnw/goplugin"
	"github.com/qxnw/goplugin/errorCode"
	"github.com/qxnw/lib4go/encoding/base64"
	"github.com/qxnw/lib4go/security/xsrf"
	"github.com/qxnw/lib4go/types"
	"github.com/qxnw/lib4go/utility"
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

	err = context.CheckInput(fields["mustFields"]...)
	if err != nil {
		status = errorCode.NOT_ACCEPTABLE
		context.Close()
		return
	}
	err = context.CheckArgs(fields["mustArgs"]...)
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

//AddAuthTimes 累加授权次数
func (c *Context) AddAuthTimes() (err error) {
	r := c.GetCookieString("_auth_times")
	times := types.ToInt(r, 0) + 1
	c.SetCookie("_auth_times", strconv.Itoa(times), 0, c.Conf.HostName, c.Param)
	if times > 3 {
		err = fmt.Errorf("授权次数超过限制次数3次:%s", c.GetUserSessionID())
		return
	}
	return
}

//ClearAuthTimes ...
func (c *Context) ClearAuthTimes() {
	c.SetCookie("_auth_times", "0", 0, c.Conf.HostName, c.Param)
}

//RedirectToWXAuth 跳转微信授权
func (c *Context) RedirectToWXAuth() {
	//记录当前地址
	request, err := c.GetHttpRequest()
	if err != nil {
		panic(err)
	}
	c.SetCookie("user_pre_url", request.RequestURI, 0, c.Conf.HostName, c.Param)
	url := c.GetAuthURL()
	c.Infof("跳转到授权地址:%s", url)
	c.Redirect(302, url)
}

//GetAuthURL 获取用户授权地址
func (c *Context) GetAuthURL() (url string) {
	return fmt.Sprintf("%s?appid=%s&redirect_uri=http://%s%s&response_type=code&scope=snsapi_base&state=#wechat:redirect",
		c.Conf.WXAuthURL,
		c.Conf.APPID,
		c.Conf.HostName,
		c.Conf.WXAuthNotifyURL)
}

//GetXSRFToken
func (c *Context) GetXSRFToken() string {
	t, _ := base64.URLDecode(utility.GetGUID())
	token := xsrf.CreateXSRFToken(c.Conf.XSRFSecret, t)
	return token
}
func (c *Context) RegisterUser(code string) (sessionid string, err error) {
	status, result, _, err := c.RPC.Request(c.Conf.AuthorizeLogin, map[string]string{
		"state": utility.GetGUID(),
		"code":  code,
	}, true)
	if err != nil || status != 200 {
		return "", fmt.Errorf("处理用户授权信息失败:%s,%s,%d:%v", c.Conf.AuthorizeLogin, code, status, err)
	}
	resultMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &resultMap)
	if err != nil {
		return "", fmt.Errorf("用户信息转换为json异常,%s,err:%v", result, err)
	}
	if sessionid, ok := resultMap["guid"]; ok {
		return sessionid.(string), nil
	}
	return "", fmt.Errorf("返回的数据中没有guid(session):%s", result)
}

//FetchUserSession 用微信code换取session
func (c *Context) FetchUserSession() (cURL string, sid string, err error) {
	code := c.GetString("code")
	if code == "" {
		err = fmt.Errorf("微信授权返回code为空")
		return
	}
	sid, err = c.RegisterUser(code)
	if err != nil || sid == "" {
		err = fmt.Errorf("获取用户session错误:sid:%s,err:%v", sid, err)
		return
	}
	cURL = c.GetCookieString("user_pre_url")
	c.SetCookie("user_pre_url", "", 0, c.Conf.HostName, c.Param)
	return
}

//SetSessionID 设置session
func (c *Context) SetSessionID(sessionID string) string {
	c.SetCookie("{@pShortName}_sid", sessionID, 3600*24*15, c.Conf.HostName, c.Param)
	return sessionID
}

//GetSDKParams 获取当session
func (c *Context) GetSDKParams() string {
	return c.GetCookieString("wx_sdk_params")
}
func (c *Context) IsLogin() (err error) {
	sid := c.GetUserSessionID()
	if sid == "" {
		err = errors.New("sid为空，判定为session过期或新用户，需要跳转授权")
		return
	}
	return
}

//GetSessionID 获取当session
func (c Context) GetUserSessionID() string {
	return c.GetCookieString("{@pShortName}_sid")
}

func (c *Context) Close() {
	contextPool.Put(c)
	if c.PluginContext != nil {
		c.PluginContext.Close()
	}
}
`
