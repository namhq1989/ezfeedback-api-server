package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s server) registerIamRoutes() {
	g := s.echo.Group("/api/iam")

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

	g.POST("/request-verification-code", func(c echo.Context) error {
		var (
			ctx = c.Get("ctx").(*appcontext.AppContext)
			req = c.Get("req").(dto.RequestVerificationCodeRequest)
			ip  = ctx.GetIP()
		)

		resp, err := s.app.RequestVerificationCode(ctx, ip, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.RequestVerificationCodeRequest](next)
	})

	g.POST("/verify-verification-code", func(c echo.Context) error {
		var (
			ctx = c.Get("ctx").(*appcontext.AppContext)
			req = c.Get("req").(dto.VerifyVerificationCodeRequest)
			ip  = ctx.GetIP()
		)

		resp, err := s.app.VerifyVerificationCode(ctx, ip, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.VerifyVerificationCodeRequest](next)
	})
}
