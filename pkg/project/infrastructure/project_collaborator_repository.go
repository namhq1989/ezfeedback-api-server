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

type ProjectCollaboratorRepository struct {
	db *database.Database
}

func NewProjectCollaboratorRepository(db *database.Database) ProjectCollaboratorRepository {
	r := ProjectCollaboratorRepository{
		db: db,
	}

	return r
}

func (r ProjectCollaboratorRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectCollaboratorRepository) getTable() *table.ProjectCollaboratorsTable {
	return table.ProjectCollaborators
}

func (r ProjectCollaboratorRepository) Create(ctx *appcontext.AppContext, collaborator domain.ProjectCollaborator) error {
	mapper := mapping.ProjectCollaboratorMapper{}
	doc, err := mapper.FromDomainToModel(collaborator)
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

func (r ProjectCollaboratorRepository) Update(ctx *appcontext.AppContext, collaborator domain.ProjectCollaborator) error {
	mapper := mapping.ProjectCollaboratorMapper{}
	doc, err := mapper.FromDomainToModel(collaborator)
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

func (r ProjectCollaboratorRepository) Delete(ctx *appcontext.AppContext, collaborator domain.ProjectCollaborator) error {
	var (
		cm = r.getTable()
	)

	stmt := cm.DELETE().WHERE(cm.ID.EQ(postgres.String(collaborator.ID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r ProjectCollaboratorRepository) FindByID(ctx *appcontext.AppContext, collaboratorID string) (*domain.ProjectCollaborator, error) {
	if !uuid.IsValidID(collaboratorID) {
		return nil, apperrors.Project.InvalidCollaborator
	}

	var c = r.getTable()

	stmt := postgres.SELECT(
		c.AllColumns,
	).
		FROM(c).
		WHERE(c.ID.EQ(postgres.String(collaboratorID)))

	var doc model.ProjectCollaborators
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectCollaboratorMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r ProjectCollaboratorRepository) FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCollaborator, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var c = r.getTable()

	stmt := postgres.SELECT(
		c.AllColumns,
	).
		FROM(c).
		WHERE(c.ProjectID.EQ(postgres.String(projectID))).
		ORDER_BY(c.CreatedAt.DESC())

	var (
		docs   = make([]model.ProjectCollaborators, 0)
		result = make([]domain.ProjectCollaborator, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectCollaboratorMapper{}
	)
	for _, doc := range docs {
		collaborator, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *collaborator)
	}
	return result, nil
}

func (r ProjectCollaboratorRepository) CountTotalByProjectID(ctx *appcontext.AppContext, projectID string) (int64, error) {
	if !uuid.IsValidID(projectID) {
		return 0, apperrors.Project.InvalidProjectID
	}

	var (
		c = r.getTable()
	)

	stmt := postgres.SELECT(
		postgres.COUNT(c.ID).AS("count_result.total"),
	).
		FROM(c).
		WHERE(
			c.ProjectID.EQ(postgres.String(projectID)),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}
