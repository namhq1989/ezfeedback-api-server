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

type ProjectRepository struct {
	db *database.Database
}

func NewProjectRepository(db *database.Database) ProjectRepository {
	r := ProjectRepository{
		db: db,
	}

	return r
}

func (r ProjectRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectRepository) getTable() *table.ProjectsTable {
	return table.Projects
}

func (r ProjectRepository) Create(ctx *appcontext.AppContext, project domain.Project) error {
	mapper := mapping.ProjectMapper{}
	doc, err := mapper.FromDomainToModel(project)
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

func (r ProjectRepository) Update(ctx *appcontext.AppContext, project domain.Project) error {
	mapper := mapping.ProjectMapper{}
	doc, err := mapper.FromDomainToModel(project)
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

func (r ProjectRepository) FindByID(ctx *appcontext.AppContext, projectID string) (*domain.Project, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var p = r.getTable()

	stmt := postgres.SELECT(
		p.AllColumns,
	).
		FROM(p).
		WHERE(p.ID.EQ(postgres.String(projectID)))

	var doc model.Projects
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r ProjectRepository) FindByUserID(ctx *appcontext.AppContext, userID string) ([]domain.Project, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	var p = r.getTable()

	stmt := postgres.SELECT(
		p.AllColumns,
	).
		FROM(p).
		WHERE(p.UserID.EQ(postgres.String(userID))).
		ORDER_BY(p.CreatedAt.DESC())

	var (
		docs   = make([]model.Projects, 0)
		result = make([]domain.Project, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectMapper{}
	)
	for _, doc := range docs {
		project, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *project)
	}
	return result, nil
}
