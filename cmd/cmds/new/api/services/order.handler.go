package handlers

var HandlerOrderTmpl = `package handlers
import (

	"{@pImportName}/libs/order"
	"github.com/qxnw/hydra/context"
)

//OrderQuery 订单查询
type OrderQuery struct {
	orderLib order.IOrderLib
	fields   map[string][]string
}

//NewOrderQuery 创建订单查询对象
func NewOrderQuery() *OrderQuery {
	return &OrderQuery{
		orderLib: &order.OrderLib{},
		fields: map[string][]string{
			"input": []string{"session_id"},
		},
	}
}

//Handle 业务处理
func (o *OrderQuery) Handle(name string, mode string, service string, ctx *context.Context) (response *context.ObjectResponse, err error) {
	response = context.GetObjectResponse()
	status, err := ctx.Input.Check(o.fields)
	if err != nil {
		response.SetError(status, err)
		return
	}

	ctx.Info("---------------查询订单信息---------------")
	sid := ctx.Input.GetString("session_id")

	ctx.Info("1. 从数据库获取订单信息")
	orderList, err := o.orderLib.QueryById(ctx, sid)
	if err != nil {
		response.SetStatus(context.ERR_UNAUTHORIZED)
		return
	}
	ctx.Info("2. 返回结果")
	response.Success(orderList)
	return
}

//Close 释放资源
func (n *OrderQuery) Close() error {
	return nil
}`
