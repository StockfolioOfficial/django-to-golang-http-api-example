package supporter

import "github.com/labstack/echo/v4"

type Route func(e *echo.Echo)

type Routes []Route

func Path(pattern string, view View) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		e.GET(pattern, view.Get)
		e.HEAD(pattern, view.Head)
		e.POST(pattern, view.Post)
		e.PUT(pattern, view.Put)
		e.PATCH(pattern, view.Patch)

		e.DELETE(pattern, view.Delete)
		e.CONNECT(pattern, view.Connect)
		e.OPTIONS(pattern, view.Options)
		e.TRACE(pattern, view.Trace)
	}
}
