package internals

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderTempl(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
