package infrastructure

import (
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
)

type NotificationReminderRepository struct {
	db                   *database.Database
	cleanupStaleDuration time.Duration
}

func NewNotificationReminderRepository(db *database.Database) NotificationReminderRepository {
	r := NotificationReminderRepository{
		db:                   db,
		cleanupStaleDuration: -1 * 24 * time.Hour,
	}

	return r
}

func (r NotificationReminderRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (NotificationReminderRepository) getTable() *table.NotificationRemindersTable {
	return table.NotificationReminders
}

func (r NotificationReminderRepository) Create(ctx *appcontext.AppContext, reminder domain.NotificationReminder) error {
	mapper := mapping.NotificationReminderMapper{}
	doc, err := mapper.FromDomainToModel(reminder)
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

func (r NotificationReminderRepository) Delete(ctx *appcontext.AppContext, reminder domain.NotificationReminder) error {
	var (
		cm = r.getTable()
	)

	stmt := cm.DELETE().WHERE(cm.ID.EQ(postgres.String(reminder.ID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r NotificationReminderRepository) FindAllExisting(ctx *appcontext.AppContext) ([]domain.NotificationReminder, error) {
	var (
		nr = r.getTable()
	)

	stmt := postgres.SELECT(
		nr.AllColumns,
	).
		FROM(nr).
		ORDER_BY(nr.CreatedAt.ASC())

	var (
		docs   = make([]model.NotificationReminders, 0)
		result = make([]domain.NotificationReminder, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.NotificationReminderMapper{}
	)
	for _, doc := range docs {
		reminder, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *reminder)
	}
	return result, nil
}

func (r NotificationReminderRepository) CleanupStale(ctx *appcontext.AppContext) error {
	var (
		nr = r.getTable()
		ts = manipulation.NowUTC().Add(r.cleanupStaleDuration)
	)

	stmt := nr.DELETE().WHERE(nr.CreatedAt.LT(postgres.TimestampzT(ts)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}
