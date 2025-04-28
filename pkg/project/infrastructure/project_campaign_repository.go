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

type ProjectCampaignRepository struct {
	db *database.Database
}

func NewProjectCampaignRepository(db *database.Database) ProjectCampaignRepository {
	r := ProjectCampaignRepository{
		db: db,
	}

	return r
}

func (r ProjectCampaignRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectCampaignRepository) getTable() *table.ProjectCampaignsTable {
	return table.ProjectCampaigns
}

func (r ProjectCampaignRepository) Create(ctx *appcontext.AppContext, campaign domain.ProjectCampaign) error {
	mapper := mapping.ProjectCampaignMapper{}
	doc, err := mapper.FromDomainToModel(campaign)
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

func (r ProjectCampaignRepository) Update(ctx *appcontext.AppContext, campaign domain.ProjectCampaign) error {
	mapper := mapping.ProjectCampaignMapper{}
	doc, err := mapper.FromDomainToModel(campaign)
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

func (r ProjectCampaignRepository) FindByID(ctx *appcontext.AppContext, campaignID string) (*domain.ProjectCampaign, error) {
	if !uuid.IsValidID(campaignID) {
		return nil, apperrors.Project.InvalidCampaign
	}

	var c = r.getTable()

	stmt := postgres.SELECT(
		c.AllColumns,
	).
		FROM(c).
		WHERE(c.ID.EQ(postgres.String(campaignID)))

	var doc model.ProjectCampaigns
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectCampaignMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r ProjectCampaignRepository) FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCampaign, error) {
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
		docs   = make([]model.ProjectCampaigns, 0)
		result = make([]domain.ProjectCampaign, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectCampaignMapper{}
	)
	for _, doc := range docs {
		campaign, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *campaign)
	}
	return result, nil
}

func (r ProjectCampaignRepository) CountTotalByProjectIDAndCampaignType(ctx *appcontext.AppContext, projectID, campaignType string) (int64, error) {
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
			c.ProjectID.EQ(postgres.String(projectID)).
				AND(c.CampaignType.EQ(postgres.String(campaignType))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}
