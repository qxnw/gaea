package order

var OrderLibTmpl = `package order
import (
	"github.com/qxnw/lib4go/db"
	"github.com/qxnw/hydra/context"
)


type IOrderLib interface {
	QueryById(context *context.Context, id string) ([]db.QueryRow, error)
}
type OrderLib struct {
}

func (o *OrderLib) QueryById(context *context.Context, id string) ([]db.QueryRow, error) {
	return context.DB.GetDataRows(sql_QUERY_ORDER, map[string]interface{}{"id": id})
}`
