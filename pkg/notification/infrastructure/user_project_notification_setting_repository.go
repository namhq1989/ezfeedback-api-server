package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type UserProjectNotificationSettingRepository struct {
	db *database.Database
}

func NewUserProjectNotificationSettingRepository(db *database.Database) UserProjectNotificationSettingRepository {
	r := UserProjectNotificationSettingRepository{
		db: db,
	}

	return r
}

func (r UserProjectNotificationSettingRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (UserProjectNotificationSettingRepository) getTable() *table.UserProjectNotificationSettingsTable {
	return table.UserProjectNotificationSettings
}

func (r UserProjectNotificationSettingRepository) Create(ctx *appcontext.AppContext, setting domain.UserProjectNotificationSetting) error {
	mapper := mapping.UserProjectNotificationSettingMapper{}
	doc, err := mapper.FromDomainToModel(setting)
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

func (r UserProjectNotificationSettingRepository) Update(ctx *appcontext.AppContext, setting domain.UserProjectNotificationSetting) error {
	mapper := mapping.UserProjectNotificationSettingMapper{}
	doc, err := mapper.FromDomainToModel(setting)
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

func (r UserProjectNotificationSettingRepository) FindByUserIDAndProjectID(ctx *appcontext.AppContext, userID, projectID string) (*domain.UserProjectNotificationSetting, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var upns = r.getTable()

	stmt := postgres.SELECT(
		upns.AllColumns,
	).
		FROM(upns).
		WHERE(upns.UserID.EQ(postgres.String(userID)).
			AND(upns.ProjectID.EQ(postgres.String(projectID))))

	var doc model.UserProjectNotificationSettings
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.UserProjectNotificationSettingMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}
