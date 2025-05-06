package httprespond

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/go-utilities/appcontext"
)

func isDataNil(data interface{}) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func sendResponse(c echo.Context, httpCode int, err error, data interface{}) error {
	var (
		ctx  = c.Get("ctx").(*appcontext.AppContext)
		lang = ctx.GetLang()
	)

	if isDataNil(data) {
		data = echo.Map{}
	}

	code, message := apperrors.GetMessage(lang.String(), err)

	return c.JSON(httpCode, echo.Map{
		"data":    data,
		"code":    code,
		"message": message,
	})
}

// R200 response success
func R200(c echo.Context, data interface{}) error {
	return sendResponse(c, http.StatusOK, apperrors.Common.Success, data)
}

// R400 bad request
func R400(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.BadRequest
	}

	// redirect to 403 if error is not_allowed
	if errors.Is(err, apperrors.Auth.NotAllowed) {
		return R403(c, err, data)
	}

	// redirect to 404 if error is not_found
	if strings.Contains(err.Error(), "not_found") {
		return R404(c, err, data)
	}

	return sendResponse(c, http.StatusBadRequest, err, data)
}

// R401 unauthorized
func R401(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.Unauthorized
	}
	return sendResponse(c, http.StatusUnauthorized, err, data)
}

// R403 forbidden
func R403(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.Forbidden
	}
	return sendResponse(c, http.StatusForbidden, err, data)
}

// R404 not found
func R404(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.NotFound
	}
	return sendResponse(c, http.StatusNotFound, err, data)
}
