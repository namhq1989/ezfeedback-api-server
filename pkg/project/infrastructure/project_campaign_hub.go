package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ProjectCampaignHub struct {
	db *database.Database
}

func NewProjectCampaignHub(db *database.Database) ProjectCampaignHub {
	r := ProjectCampaignHub{
		db: db,
	}

	return r
}

func (r ProjectCampaignHub) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectCampaignHub) getTable() *table.ProjectCampaignsTable {
	return table.ProjectCampaigns
}

func (r ProjectCampaignHub) FindProjectCampaignByID(ctx *appcontext.AppContext, id string) (*domain.ProjectCampaignHubData, error) {
	var (
		pc = r.getTable().AS("pc")
		p  = table.Projects.AS("p")
		ps = table.ProjectSettings.AS("ps")
	)

	stmt := postgres.SELECT(
		pc.ID, pc.Name, pc.CampaignType, pc.Status,
		p.ID, p.Title, p.Status,
		ps.Domain,
	).
		FROM(
			pc.LEFT_JOIN(p, pc.ProjectID.EQ(p.ID)).
				LEFT_JOIN(ps, p.ID.EQ(ps.ProjectID)),
		).
		WHERE(pc.ID.EQ(postgres.String(id)))

	var doc mapping.ProjectCampaignHubData
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectCampaignHubDataMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}
