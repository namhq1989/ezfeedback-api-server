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
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure/mapping"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure/template"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type NotificationRepository struct {
	db                   *database.Database
	cleanupStaleDuration time.Duration
}

func NewNotificationRepository(db *database.Database) NotificationRepository {
	r := NotificationRepository{
		db:                   db,
		cleanupStaleDuration: -1 * 6 * 30 * 24 * time.Hour,
	}

	return r
}

func (r NotificationRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (NotificationRepository) getTable() *table.NotificationsTable {
	return table.Notifications
}

func (r NotificationRepository) Create(ctx *appcontext.AppContext, notification domain.Notification) error {
	mapper := mapping.NotificationMapper{}
	doc, err := mapper.FromDomainToModel(notification)
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

func (r NotificationRepository) Update(ctx *appcontext.AppContext, notification domain.Notification) error {
	mapper := mapping.NotificationMapper{}
	doc, err := mapper.FromDomainToModel(notification)
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

func (r NotificationRepository) Delete(ctx *appcontext.AppContext, notification domain.Notification) error {
	var (
		n = r.getTable()
	)

	stmt := n.DELETE().WHERE(n.ID.EQ(postgres.String(notification.ID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r NotificationRepository) FindByID(ctx *appcontext.AppContext, id string) (*domain.Notification, error) {
	if !uuid.IsValidID(id) {
		return nil, apperrors.Notification.InvalidID
	}

	var n = r.getTable()

	stmt := postgres.SELECT(
		n.AllColumns,
	).
		FROM(n).
		WHERE(n.ID.EQ(postgres.String(id)))

	var doc model.Notifications
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.NotificationMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r NotificationRepository) FindWithFilter(ctx *appcontext.AppContext, filter domain.NotificationFilter) ([]domain.Notification, error) {
	var (
		n         = r.getTable()
		offset    = filter.Limit * filter.Page
		whereStmt = n.UserID.EQ(postgres.String(filter.UserID))
	)

	stmt := postgres.SELECT(
		n.ID, n.UserID, n.Type, n.IsRead, n.Metadata, n.CreatedAt,
	).
		FROM(n).
		WHERE(whereStmt).
		LIMIT(filter.Limit).
		OFFSET(offset).
		ORDER_BY(n.CreatedAt.DESC())

	var (
		docs   = make([]model.Notifications, 0)
		result = make([]domain.Notification, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.NotificationMapper{}
	)
	for _, doc := range docs {
		notification, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *notification)
	}
	return result, nil
}

func (r NotificationRepository) CountWithFilter(ctx *appcontext.AppContext, filter domain.NotificationFilter) (int64, error) {
	var (
		n         = r.getTable()
		whereStmt = n.UserID.EQ(postgres.String(filter.UserID))
	)

	stmt := postgres.SELECT(
		postgres.COUNT(n.ID).AS("count_result.total"),
	).
		FROM(n).
		WHERE(whereStmt)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}

func (r NotificationRepository) CleanupStale(ctx *appcontext.AppContext) error {
	var (
		n  = r.getTable()
		ts = manipulation.NowUTC().Add(r.cleanupStaleDuration)
	)

	stmt := n.DELETE().WHERE(n.CreatedAt.LT(postgres.TimestampzT(ts)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (NotificationRepository) GenerateNewFeedbackContent(_ *appcontext.AppContext, language string, projectTitle string) string {
	return template.NewFeedbackContent(language, projectTitle)
}
