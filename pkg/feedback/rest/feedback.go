package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s server) registerFeedbackRoutes() {
	g := s.echo.Group("/api/feedback")

	g.GET("/ping", func(c echo.Context) error {
		var (
			ctx = c.Get("ctx").(*appcontext.AppContext)
			req = c.Get("req").(dto.PingRequest)
		)

		resp, err := s.app.Ping(ctx, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.PingRequest](next)
	})

	g.POST("", func(c echo.Context) error {
		var (
			ctx    = c.Get("ctx").(*appcontext.AppContext)
			req    = c.Get("req").(dto.CreateFeedbackRequest)
			ip     = ctx.GetIP()
			origin = c.Request().Header.Get("Origin")
		)

		resp, err := s.app.CreateFeedback(ctx, ip, origin, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateFeedbackRequest](next)
	})
}
