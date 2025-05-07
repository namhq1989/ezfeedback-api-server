package rest

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
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

	g.GET("", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetFeedbacksRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetFeedbacks(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetFeedbacksRequest](next)
	})

	g.GET("/count", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.CountFeedbacksRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.CountFeedbacks(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CountFeedbacksRequest](next)
	})

	g.POST("", func(c echo.Context) error {
		var (
			ctx    = c.Get("ctx").(*appcontext.AppContext)
			req    = c.Get("req").(dto.CreateFeedbackRequest)
			ip     = ctx.GetIP()
			origin = c.Request().Header.Get("Origin")
		)

		domainName, err := manipulation.GetRootDomain(origin)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		resp, err := s.app.CreateFeedback(ctx, ip, domainName, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, createFeedbackRateLimiter(), func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateFeedbackRequest](next)
	})

	g.PATCH("/:id/state", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ChangeFeedbackStateRequest)
			performerID = ctx.GetUserID()
			feedbackID  = c.Param("id")
		)

		resp, err := s.app.ChangeFeedbackState(ctx, performerID, feedbackID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ChangeFeedbackStateRequest](next)
	})

	g.POST("/:id/reply", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.CreateFeedbackReplyRequest)
			performerID = ctx.GetUserID()
			feedbackID  = c.Param("id")
		)

		resp, err := s.app.CreateFeedbackReply(ctx, performerID, feedbackID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.CreateFeedbackReplyRequest](next)
	})

	g.PUT("/:feedbackId/reply/:replyId", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.UpdateFeedbackReplyRequest)
			performerID = ctx.GetUserID()
			feedbackID  = c.Param("feedbackId")
			replyID     = c.Param("replyId")
		)

		resp, err := s.app.UpdateFeedbackReply(ctx, performerID, feedbackID, replyID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireSignedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.UpdateFeedbackReplyRequest](next)
	})
}

func createFeedbackRateLimiter() echo.MiddlewareFunc {
	feedbackLimiterStore := middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      0.1, // 0.1 requests per second = 1 request per 10 seconds
			Burst:     1,   // Only allow 1 request at a time
			ExpiresIn: 1 * time.Minute,
		},
	)

	feedbackLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store:   feedbackLimiterStore,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, _ error) error {
			return context.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Rate limiting error occurred",
			})
		},
		DenyHandler: func(context echo.Context, _ string, _ error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Your feedback is important to us. Please wait a moment before submitting another response",
			})
		},
	}

	return middleware.RateLimiterWithConfig(feedbackLimiterConfig)
}
