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

type ProjectCategoryRepository struct {
	db *database.Database
}

func NewProjectCategoryRepository(db *database.Database) ProjectCategoryRepository {
	r := ProjectCategoryRepository{
		db: db,
	}

	return r
}

func (r ProjectCategoryRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectCategoryRepository) getTable() *table.ProjectCategoriesTable {
	return table.ProjectCategories
}

func (r ProjectCategoryRepository) Create(ctx *appcontext.AppContext, category domain.ProjectCategory) error {
	mapper := mapping.ProjectCategoryMapper{}
	doc, err := mapper.FromDomainToModel(category)
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

func (r ProjectCategoryRepository) Update(ctx *appcontext.AppContext, category domain.ProjectCategory) error {
	mapper := mapping.ProjectCategoryMapper{}
	doc, err := mapper.FromDomainToModel(category)
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

func (r ProjectCategoryRepository) FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCategory, error) {
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
		docs   = make([]model.ProjectCategories, 0)
		result = make([]domain.ProjectCategory, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectCategoryMapper{}
	)
	for _, doc := range docs {
		category, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *category)
	}
	return result, nil
}

func (r ProjectCategoryRepository) CountTotalByProjectID(ctx *appcontext.AppContext, projectID string) (int64, error) {
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
