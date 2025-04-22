package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectSettingRepository struct {
	db *database.Database
}

func NewProjectSettingRepository(db *database.Database) ProjectSettingRepository {
	r := ProjectSettingRepository{
		db: db,
	}

	return r
}

func (r ProjectSettingRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectSettingRepository) getTable() *table.ProjectSettingsTable {
	return table.ProjectSettings
}

func (r ProjectSettingRepository) Create(ctx *appcontext.AppContext, setting domain.ProjectSetting) error {
	mapper := mapping.ProjectSettingMapper{}
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

func (r ProjectSettingRepository) Update(ctx *appcontext.AppContext, setting domain.ProjectSetting) error {
	mapper := mapping.ProjectSettingMapper{}
	doc, err := mapper.FromDomainToModel(setting)
	if err != nil {
		return err
	}

	stmt := r.getTable().UPDATE(
		r.getTable().AllColumns,
	).
		MODEL(doc).
		WHERE(
			r.getTable().ProjectID.EQ(postgres.String(doc.ID)),
		)

	_, err = stmt.ExecContext(ctx.Context(), r.getDB())
	if err != nil {
		if isDuplicated, duplicateErr := r.db.IsDuplicatedError(err); isDuplicated {
			err = duplicateErr
		}
	}
	return err
}

func (r ProjectSettingRepository) FindByProjectID(ctx *appcontext.AppContext, projectID string) (*domain.ProjectSetting, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var s = r.getTable()

	stmt := postgres.SELECT(
		s.AllColumns,
	).
		FROM(s).
		WHERE(s.ProjectID.EQ(postgres.String(projectID)))

	var doc model.ProjectSettings
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectSettingMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}
