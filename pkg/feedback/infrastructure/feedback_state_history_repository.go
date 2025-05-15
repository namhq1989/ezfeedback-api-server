package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
)

type FeedbackStateHistoryRepository struct {
	db *database.Database
}

func NewFeedbackStateHistoryRepository(db *database.Database) FeedbackStateHistoryRepository {
	r := FeedbackStateHistoryRepository{
		db: db,
	}

	return r
}

func (r FeedbackStateHistoryRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (FeedbackStateHistoryRepository) getTable() *table.FeedbackStateHistoriesTable {
	return table.FeedbackStateHistories
}

func (r FeedbackStateHistoryRepository) Create(ctx *appcontext.AppContext, history domain.FeedbackStateHistory) error {
	mapper := mapping.FeedbackStateHistoryMapper{}
	doc, err := mapper.FromDomainToModel(history)
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

func (r FeedbackStateHistoryRepository) FindWithFilter(ctx *appcontext.AppContext, filter domain.FeedbackStateHistoryFilter) ([]domain.FeedbackStateHistory, error) {
	var (
		fsh       = r.getTable()
		offset    = filter.Limit * filter.Page
		whereStmt = fsh.FeedbackID.EQ(postgres.String(filter.FeedbackID))
	)

	stmt := postgres.SELECT(
		fsh.ID, fsh.UserID, fsh.State, fsh.CreatedAt,
	).
		FROM(fsh).
		WHERE(whereStmt).
		LIMIT(filter.Limit).
		OFFSET(offset).
		ORDER_BY(fsh.CreatedAt.DESC())

	var (
		docs   = make([]model.FeedbackStateHistories, 0)
		result = make([]domain.FeedbackStateHistory, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.FeedbackStateHistoryMapper{}
	)
	for _, doc := range docs {
		history, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *history)
	}
	return result, nil
}
