package handlers

var HandlerOrderTmpl = `package services
import (
	"github.com/qxnw/hydra/context"
)


type MyService struct {
	fields   map[string][]string
}

func NewMyService() *MyService {
	return &MyService{
		fields: map[string][]string{
			"input": []string{"id"},
		},
	}
}

//Handle 业务处理
func (o *MyService) Handle(name string, mode string, service string, ctx *context.Context) (response *context.ObjectResponse, err error) {
	response = context.GetObjectResponse()
	status, err := ctx.Input.Check(o.fields)
	if err != nil {
		response.SetError(status, err)
		return
	}

	ctx.Info("---------------业务处理---------------")
	id := ctx.Input.GetString("id")
	ctx.Info("2. 返回结果")
	response.Success(map[string]string{
		"id":id,
	})
	return
}


//Close 释放资源
func (n *MyService) Close() error {
	return nil
}`
