package handlers

var HandlerOrderTmpl = `package handlers
import (
	"{@pImportName}/context"

	"{@pImportName}/libs/order"

	"github.com/qxnw/goplugin"
	"github.com/qxnw/goplugin/errorCode"
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
func (o *OrderQuery) Handle(service string, ctx goplugin.Context, rpc goplugin.RPCInvoker) (status int, result interface{}, p map[string]interface{}, err error) {
	context, status, p, err := context.GetContext(ctx, rpc, o.fields)
	if err != nil {
		return
	}

	context.Info("---------------查询订单信息---------------")
	defer context.Close()
	sid, _ := context.Input.Get("session_id")

	context.Info("1. 从数据库获取订单信息")
	orderList, err := o.orderLib.QueryById(context, sid)
	if err != nil {
		return errorCode.UNAUTHORIZED, nil, p, err
	}
	context.Info("2. 返回结果")
	return errorCode.OK, orderList, p, nil
}

//Close 释放资源
func (n *OrderQuery) Close() error {
	return nil
}
`
