package supporter

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type View interface {
	Get(ctx echo.Context) error
	Head(ctx echo.Context) error
	Post(ctx echo.Context) error
	Put(ctx echo.Context) error
	Patch(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Connect(ctx echo.Context) error
	Options(ctx echo.Context) error
	Trace(ctx echo.Context) error
	impl()
}

type APIView struct { }

func (APIView) Get(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Head(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Post(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Put(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Patch(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Delete(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Connect(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Options(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) Trace(ctx echo.Context) error {
	return ctx.NoContent(http.StatusMethodNotAllowed)
}

func (APIView) impl() {}