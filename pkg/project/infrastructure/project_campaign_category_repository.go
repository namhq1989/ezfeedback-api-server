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

type ProjectCampaignCategoryRepository struct {
	db *database.Database
}

func NewProjectCampaignCategoryRepository(db *database.Database) ProjectCampaignCategoryRepository {
	r := ProjectCampaignCategoryRepository{
		db: db,
	}

	return r
}

func (r ProjectCampaignCategoryRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectCampaignCategoryRepository) getTable() *table.ProjectCampaignCategoriesTable {
	return table.ProjectCampaignCategories
}

func (r ProjectCampaignCategoryRepository) Create(ctx *appcontext.AppContext, campaignCategory domain.ProjectCampaignCategory) error {
	mapper := mapping.ProjectCampaignCategoryMapper{}
	doc, err := mapper.FromDomainToModel(campaignCategory)
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

func (r ProjectCampaignCategoryRepository) Delete(ctx *appcontext.AppContext, campaignCategory domain.ProjectCampaignCategory) error {
	var (
		cm = r.getTable()
	)

	stmt := cm.DELETE().WHERE(cm.ID.EQ(postgres.String(campaignCategory.ID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r ProjectCampaignCategoryRepository) FindByID(ctx *appcontext.AppContext, campaignCategoryID string) (*domain.ProjectCampaignCategory, error) {
	if !uuid.IsValidID(campaignCategoryID) {
		return nil, apperrors.Common.InvalidID
	}

	var cm = r.getTable()

	stmt := postgres.SELECT(
		cm.AllColumns,
	).
		FROM(cm).
		WHERE(cm.ID.EQ(postgres.String(campaignCategoryID)))

	var doc model.ProjectCampaignCategories
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectCampaignCategoryMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r ProjectCampaignCategoryRepository) FindByCampaignID(ctx *appcontext.AppContext, campaignID string) ([]domain.ProjectCampaignCategory, error) {
	if !uuid.IsValidID(campaignID) {
		return nil, apperrors.Project.InvalidCampaign
	}

	var cm = r.getTable()

	stmt := postgres.SELECT(
		cm.AllColumns,
	).
		FROM(cm).
		WHERE(cm.CampaignID.EQ(postgres.String(campaignID))).
		ORDER_BY(cm.CreatedAt.DESC())

	var (
		docs   = make([]model.ProjectCampaignCategories, 0)
		result = make([]domain.ProjectCampaignCategory, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectCampaignCategoryMapper{}
	)
	for _, doc := range docs {
		campaignCategory, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *campaignCategory)
	}
	return result, nil
}
