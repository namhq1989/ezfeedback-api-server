package infrastructure

import (
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackHub struct {
	db *database.Database
}

func NewFeedbackHub(db *database.Database) FeedbackHub {
	r := FeedbackHub{
		db: db,
	}

	return r
}

func (r FeedbackHub) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (FeedbackHub) getTable() *table.FeedbacksTable {
	return table.Feedbacks
}

func (r FeedbackHub) CountFeedbackForProjectSinceTimestamp(ctx *appcontext.AppContext, projectID string, timestamp time.Time) (int64, error) {
	if !uuid.IsValidID(projectID) {
		return 0, apperrors.Project.InvalidProjectID
	}

	var (
		f = r.getTable()
	)

	stmt := postgres.SELECT(
		postgres.COUNT(f.ID).AS("count_result.total"),
	).
		FROM(f).
		WHERE(
			f.ProjectID.EQ(postgres.String(projectID)).
				AND(f.CreatedAt.GT_EQ(postgres.TimestampzT(timestamp))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}

func (r FeedbackHub) FindFeedbackForProjectSinceTimestamp(ctx *appcontext.AppContext, projectID string, timestamp time.Time, limit int64) ([]domain.Feedback, error) {
	if !uuid.IsValidID(projectID) {
		return make([]domain.Feedback, 0), apperrors.Project.InvalidProjectID
	}

	var (
		f = r.getTable()
	)

	stmt := postgres.SELECT(
		f.ID, f.ProjectID, f.AppUserID, f.Email, f.Content,
		f.Rating, f.CampaignType, f.CreatedAt,
	).
		FROM(f).
		WHERE(
			f.ProjectID.EQ(postgres.String(projectID)).
				AND(f.CreatedAt.GT_EQ(postgres.TimestampzT(timestamp))),
		).
		ORDER_BY(f.CreatedAt.DESC()).
		LIMIT(limit)

	var (
		docs   = make([]model.Feedbacks, 0)
		result = make([]domain.Feedback, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.FeedbackMapper{}
	)
	for _, doc := range docs {
		feedback, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *feedback)
	}
	return result, nil
}
