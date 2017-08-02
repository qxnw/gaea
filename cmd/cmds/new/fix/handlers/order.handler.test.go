package handlers

var HandlerTestTmpl = `package handlers
import (
	"testing"
	"{@pImportName}/context"

	"github.com/qxnw/goplugin"
	"github.com/qxnw/goplugin/errorCode"
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
	context := goplugin.NewTContext()
	rpc := &goplugin.TRPC{}
	s, r, p, err := query.Handle("", context, rpc)
	ut.Expect(t, s, errorCode.NOT_ACCEPTABLE)
	ut.Expect(t, r, nil)
	ut.Expect(t, len(p), 0)
	ut.Refute(t, err, "")
}

func TestHandle11(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	context := goplugin.NewTContext()
	context.Input.Set("session_Id", "1")
	rpc := &goplugin.TRPC{}
	s, r, p, err := query.Handle("", context, rpc)
	ut.Expect(t, s, errorCode.NOT_ACCEPTABLE)
	ut.Expect(t, r, nil)
	ut.Expect(t, len(p), 0)
	ut.Refute(t, err, "")
}
func TestHandle2(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	context := goplugin.NewTContext()
	context.Input.Set("session_id", "123445")
	rpc := &goplugin.TRPC{}
	s, r, p, err := query.Handle("", context, rpc)
	ut.Expect(t, s, errorCode.OK)
	ut.Expect(t, len(r.([]db.QueryRow)), 0)
	ut.Expect(t, len(p), 0)
	ut.Refute(t, err, "")
}
func TestHandle3(t *testing.T) {
	query := NewOrderQuery()
	query.orderLib = &orderMock{}
	context := goplugin.NewTContext()
	context.Input.Set("session_id", "1")
	rpc := &goplugin.TRPC{}
	s, r, p, err := query.Handle("", context, rpc)
	ut.Expect(t, s, errorCode.OK)
	ut.Expect(t, len(r.([]db.QueryRow)), 1)
	ut.Expect(t, len(p), 0)
	ut.Refute(t, err, "")
}
`
