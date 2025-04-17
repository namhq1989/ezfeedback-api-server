package validation

import (
	"net/http"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
)

func ValidateHTTPPayload[T any](next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			req T
		)

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"data":    nil,
				"code":    "bad_request",
				"message": "Bad request",
			})
		}

		if v := validate.Struct(req); !v.Validate() {
			return httprespond.R400(c, v.Errors.OneError(), nil)
		}

		// assign to context
		c.Set("req", req)
		return next(c)
	}
}
