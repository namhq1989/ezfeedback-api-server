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
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectInvitationRepository struct {
	db                   *database.Database
	cleanupStaleDuration time.Duration
}

func NewProjectInvitationRepository(db *database.Database) ProjectInvitationRepository {
	r := ProjectInvitationRepository{
		db:                   db,
		cleanupStaleDuration: -1 * 15 * 24 * time.Hour,
	}

	return r
}

func (r ProjectInvitationRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (ProjectInvitationRepository) getTable() *table.ProjectInvitationsTable {
	return table.ProjectInvitations
}

func (r ProjectInvitationRepository) Create(ctx *appcontext.AppContext, invitation domain.ProjectInvitation) error {
	mapper := mapping.ProjectInvitationMapper{}
	doc, err := mapper.FromDomainToModel(invitation)
	if err != nil {
		return err
	}

	stmt := r.getTable().INSERT(
		r.getTable().AllColumns,
	).
		MODEL(doc)

	_, err = stmt.ExecContext(ctx.Context(), r.getDB())
	if err != nil {
		if isDuplicated, duplicateErr := r.db.IsDuplicatedError(err); isDuplicated {
			err = duplicateErr
		}
	}
	return err
}

func (r ProjectInvitationRepository) Update(ctx *appcontext.AppContext, invitation domain.ProjectInvitation) error {
	mapper := mapping.ProjectInvitationMapper{}
	doc, err := mapper.FromDomainToModel(invitation)
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
	return err
}

func (r ProjectInvitationRepository) Delete(ctx *appcontext.AppContext, invitationID string) error {
	if !uuid.IsValidID(invitationID) {
		return apperrors.Project.InvalidInvitation
	}

	var pi = r.getTable()

	stmt := pi.DELETE().WHERE(pi.ID.EQ(postgres.String(invitationID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r ProjectInvitationRepository) FindByID(ctx *appcontext.AppContext, invitationID string) (*domain.ProjectInvitation, error) {
	if !uuid.IsValidID(invitationID) {
		return nil, apperrors.Project.InvalidInvitation
	}

	var pi = r.getTable()

	stmt := postgres.SELECT(
		pi.AllColumns,
	).
		FROM(pi).
		WHERE(pi.ID.EQ(postgres.String(invitationID)))

	var doc model.ProjectInvitations
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.ProjectInvitationMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r ProjectInvitationRepository) FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectInvitation, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var pi = r.getTable()

	stmt := postgres.SELECT(
		pi.AllColumns,
	).
		FROM(pi).
		WHERE(pi.ProjectID.EQ(postgres.String(projectID))).
		ORDER_BY(pi.CreatedAt.DESC())

	var (
		docs   = make([]model.ProjectInvitations, 0)
		result = make([]domain.ProjectInvitation, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectInvitationMapper{}
	)
	for _, doc := range docs {
		invitation, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *invitation)
	}
	return result, nil
}

func (r ProjectInvitationRepository) FindByEmail(ctx *appcontext.AppContext, email string) ([]domain.ProjectInvitation, error) {
	var pi = r.getTable()

	stmt := postgres.SELECT(
		pi.AllColumns,
	).
		FROM(pi).
		WHERE(
			pi.Email.EQ(postgres.String(email)).
				AND(pi.Status.EQ(postgres.String(domain.ProjectInvitationStatusPending.String()))),
		).
		ORDER_BY(pi.CreatedAt.DESC())

	var (
		docs   = make([]model.ProjectInvitations, 0)
		result = make([]domain.ProjectInvitation, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.ProjectInvitationMapper{}
	)
	for _, doc := range docs {
		invitation, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *invitation)
	}
	return result, nil
}

func (r ProjectInvitationRepository) MarkExpiredInvitations(ctx *appcontext.AppContext) error {
	var (
		pi  = r.getTable()
		now = manipulation.NowUTC()
	)

	stmt := pi.UPDATE(pi.Status, pi.UpdatedAt).
		SET(
			postgres.String(domain.ProjectInvitationStatusExpired.String()),
			postgres.TimestampzT(now),
		).
		WHERE(
			pi.ExpiresAt.LT(postgres.TimestampzT(now)).
				AND(pi.Status.EQ(postgres.String(domain.ProjectInvitationStatusPending.String()))),
		)

	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r ProjectInvitationRepository) CleanupStale(ctx *appcontext.AppContext) error {
	var (
		pi = r.getTable()
		ts = manipulation.NowUTC().Add(r.cleanupStaleDuration)
	)

	stmt := pi.DELETE().WHERE(
		pi.Status.EQ(postgres.String(domain.ProjectInvitationStatusExpired.String())).
			AND(pi.UpdatedAt.LT(postgres.TimestampzT(ts))),
	)
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r ProjectInvitationRepository) CountPendingByProjectID(ctx *appcontext.AppContext, projectID string) (int64, error) {
	if !uuid.IsValidID(projectID) {
		return 0, apperrors.Project.InvalidProjectID
	}

	var (
		pi = r.getTable()
	)

	stmt := postgres.SELECT(
		postgres.COUNT(pi.ID).AS("count_result.total"),
	).
		FROM(pi).
		WHERE(
			pi.ProjectID.EQ(postgres.String(projectID)).
				AND(pi.Status.EQ(postgres.String(domain.ProjectInvitationStatusPending.String()))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}
