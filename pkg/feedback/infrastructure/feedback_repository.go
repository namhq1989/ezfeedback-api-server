package infrastructure

import (
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackRepository struct {
	db *database.Database
}

func NewFeedbackRepository(db *database.Database) FeedbackRepository {
	r := FeedbackRepository{
		db: db,
	}

	return r
}

func (r FeedbackRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (FeedbackRepository) getTable() *table.FeedbacksTable {
	return table.Feedbacks
}

func (r FeedbackRepository) Create(ctx *appcontext.AppContext, feedback domain.Feedback) error {
	mapper := mapping.FeedbackMapper{}
	doc, err := mapper.FromDomainToModel(feedback)
	if err != nil {
		return err
	}

	stmt := r.getTable().INSERT(
		r.getTable().AllColumns,
	).
		MODEL(doc)

	_, err = stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r FeedbackRepository) Update(ctx *appcontext.AppContext, feedback domain.Feedback) error {
	mapper := mapping.FeedbackMapper{}
	doc, err := mapper.FromDomainToModel(feedback)
	if err != nil {
		return err
	}

	stmt := r.getTable().UPDATE(
		r.getTable().AllColumns,
	).
		MODEL(doc).
		WHERE(
			r.getTable().ID.EQ(postgres.String(doc.ID)),
		)

	_, err = stmt.ExecContext(ctx.Context(), r.getDB())
	if err != nil {
		if isDuplicated, duplicateErr := r.db.IsDuplicatedError(err); isDuplicated {
			err = duplicateErr
		}
	}
	return err
}

func (r FeedbackRepository) FindWithFilter(ctx *appcontext.AppContext, filter domain.FeedbackFilter) ([]domain.Feedback, error) {
	var (
		f         = r.getTable()
		whereStmt = f.ProjectID.EQ(postgres.String(filter.ProjectID))
	)

	if filter.CampaignID != "" {
		whereStmt = whereStmt.AND(f.CampaignID.EQ(postgres.String(filter.CampaignID)))
	}

	if filter.Keyword != "" {
		whereStmt = whereStmt.AND(postgres.RawBool("s.search_vector @@ to_tsquery($keyword)", postgres.RawArgs{
			"$keyword": filter.Keyword,
		}))
	}

	if filter.CategoryID != "" {
		whereStmt = whereStmt.AND(f.CategoryID.EQ(postgres.String(filter.CategoryID)))
	}

	if filter.Rating != 0 {
		whereStmt = whereStmt.AND(f.Rating.EQ(postgres.Int32(filter.Rating)))
	}

	if filter.State.IsValid() {
		whereStmt = whereStmt.AND(f.State.EQ(postgres.NewEnumValue(filter.State.String())))
	}

	stmt := postgres.SELECT(
		f.ID, f.CampaignID, f.AppUserID, f.Email, f.CategoryID, f.Content,
		f.Rating, f.IsAnonymous, f.State, f.CampaignType, f.IP, f.CountryCode, f.CreatedAt,
	).
		FROM(f).
		WHERE(whereStmt).
		ORDER_BY(f.CreatedAt.DESC())

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

func (r FeedbackRepository) CountMonthlyUsageForProject(ctx *appcontext.AppContext, projectID string) (int64, error) {
	if !uuid.IsValidID(projectID) {
		return 0, apperrors.Project.InvalidProjectID
	}

	var (
		f            = r.getTable()
		now          = manipulation.NowUTC()
		startOfMonth = manipulation.StartOfMonth(now)
	)

	stmt := postgres.SELECT(
		postgres.COUNT(f.ID).AS("count_result.total"),
	).
		FROM(f).
		WHERE(
			f.ProjectID.EQ(postgres.String(projectID)).
				AND(f.CreatedAt.GT_EQ(postgres.TimestampzT(startOfMonth))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}

func (r FeedbackRepository) CountProjectTotalCreatedTodayByIp(ctx *appcontext.AppContext, projectID, ip string) (int64, error) {
	if !uuid.IsValidID(projectID) {
		return 0, apperrors.Project.InvalidProjectID
	}

	var (
		f            = r.getTable()
		startOfToday = manipulation.StartOfDay(manipulation.NowUTC())
	)

	stmt := postgres.SELECT(
		postgres.COUNT(f.ID).AS("count_result.total"),
	).
		FROM(f).
		WHERE(
			f.ProjectID.EQ(postgres.String(projectID)).
				AND(f.IP.EQ(postgres.String(ip)).
					AND(f.CreatedAt.GT_EQ(postgres.TimestampzT(startOfToday)))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}

func (r FeedbackRepository) CountFeedbackForProjectSinceTimestamp(ctx *appcontext.AppContext, projectID string, timestamp time.Time) (int64, error) {
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

func (r FeedbackRepository) FindFeedbackForProjectSinceTimestamp(ctx *appcontext.AppContext, projectID string, timestamp time.Time, limit int64) ([]domain.Feedback, error) {
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
