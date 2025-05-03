package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s server) registerNotificationRoutes() {
	g := s.echo.Group("/api/notification")

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

	g.GET("/project-setting", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetUserProjectNotificationSettingRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetUserProjectNotificationSetting(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetUserProjectNotificationSettingRequest](next)
	})

	g.PUT("/project-setting", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.UpdateUserProjectNotificationSettingRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.UpdateUserProjectNotificationSetting(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.UpdateUserProjectNotificationSettingRequest](next)
	})
}
