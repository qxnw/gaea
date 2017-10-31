package handlers

var HandlerTestTmpl = `
package handlers
import (
	"testing"

	"github.com/qxnw/hydra/context"

	"github.com/qxnw/lib4go/db"
	"github.com/qxnw/lib4go/ut"
)

type orderMock struct {
}

func (o *orderMock) QueryById(context *context.Context, id string) ([]db.QueryRow, error) {
	if id == "1" {
		return make([]db.QueryRow, 1, 1), nil
	}
	return nil, nil
}

func TestHandle1(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	ctx := context.NewTContext(nil)
	r, err := query.Handle("", "", "", ctx)
	ut.Refute(t, err, nil)
	ut.Expect(t, r.GetStatus(err), context.ERR_NOT_ACCEPTABLE)
}

func TestHandle11(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	ctx := context.NewTContext(nil)
	ctx.Input.Input.Set("session_Id", "1")
	r, err := query.Handle("", "", "", ctx)
	ut.Refute(t, err, nil)
	ut.Expect(t, r.GetStatus(err), context.ERR_NOT_ACCEPTABLE)
}
func TestHandle2(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	ctx := context.NewTContext(nil)
	ctx.Input.Input.Set("session_id", "123")
	r, err := query.Handle("", "", "", ctx)
	ut.Expect(t, r.GetStatus(err), context.ERR_OK)
	ut.Expect(t, len(r.GetContent(nil).([]db.QueryRow)), 0)
	ut.Expect(t, len(r.GetParams()), 0)
	ut.Expect(t, err, nil)
}
func TestHandle3(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	ctx := context.NewTContext(nil)
	ctx.Input.Input.Set("session_id", "1")
	r, err := query.Handle("", "", "", ctx)
	ut.Expect(t, r.GetStatus(err), context.ERR_OK)
	ut.Expect(t, len(r.GetContent(err).([]db.QueryRow)), 1)
	ut.Expect(t, len(r.GetParams()), 0)
	ut.Expect(t, err, nil)
}`
