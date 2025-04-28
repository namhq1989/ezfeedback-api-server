package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s server) registerProjectRoutes() {
	g := s.echo.Group("/api/project")

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

	g.GET("", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetProjectsRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetProjects(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetProjectsRequest](next)
	})

	g.POST("", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.CreateProjectRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.CreateProject(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateProjectRequest](next)
	})

	g.PUT("/:id", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.UpdateProjectRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("id")
		)

		resp, err := s.app.UpdateProject(ctx, performerID, projectID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.UpdateProjectRequest](next)
	})

	g.PATCH("/:id/status", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ChangeProjectStatusRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("id")
		)

		resp, err := s.app.ChangeProjectStatus(ctx, performerID, projectID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ChangeProjectStatusRequest](next)
	})

	g.POST("/:id/category", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.CreateProjectCategoryRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("id")
		)

		resp, err := s.app.CreateProjectCategory(ctx, performerID, projectID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateProjectCategoryRequest](next)
	})

	g.PUT("/:projectId/category/:categoryId", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.UpdateProjectCategoryRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("projectId")
			categoryID  = c.Param("categoryId")
		)

		resp, err := s.app.UpdateProjectCategory(ctx, performerID, projectID, categoryID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.UpdateProjectCategoryRequest](next)
	})

	g.PATCH("/:projectId/category/:categoryId/status", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ChangeProjectCategoryStatusRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("projectId")
			categoryID  = c.Param("categoryId")
		)

		resp, err := s.app.ChangeProjectCategoryStatus(ctx, performerID, projectID, categoryID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ChangeProjectCategoryStatusRequest](next)
	})

	g.POST("/:id/campaign", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.CreateProjectCampaignRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("id")
		)

		resp, err := s.app.CreateProjectCampaign(ctx, performerID, projectID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateProjectCampaignRequest](next)
	})

	g.PUT("/:projectId/campaign/:campaignId", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.UpdateProjectCampaignRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("projectId")
			campaignID  = c.Param("campaignId")
		)

		resp, err := s.app.UpdateProjectCampaign(ctx, performerID, projectID, campaignID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.UpdateProjectCampaignRequest](next)
	})

	g.PATCH("/:projectId/campaign/:campaignId/status", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ChangeProjectCampaignStatusRequest)
			performerID = ctx.GetUserID()
			projectID   = c.Param("projectId")
			campaignID  = c.Param("campaignId")
		)

		resp, err := s.app.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ChangeProjectCampaignStatusRequest](next)
	})
}
