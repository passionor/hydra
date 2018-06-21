package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type QueryHandler struct {
	container component.IContainer
}

func NewQueryHandler(container component.IContainer) (u *QueryHandler) {
	return &QueryHandler{container: container}
}

func (u *QueryHandler) GetHandle(ctx *context.Context) (r interface{}) {
	tp := ctx.Request.GetInt("t", context.CT_JSON)
	ctx.Response.SetContentType(context.ContentTypes[tp])

	return map[string]interface{}{
		"a": "b",
	}
}
func (u *QueryHandler) Handle(ctx *context.Context) (r interface{}) {
	return "success"
}
