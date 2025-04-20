package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
)

type VerificationCodeRepository struct {
	db *database.Database
}

func NewVerificationCodeRepository(db *database.Database) VerificationCodeRepository {
	r := VerificationCodeRepository{
		db: db,
	}

	return r
}

func (r VerificationCodeRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (VerificationCodeRepository) getTable() *table.VerificationCodesTable {
	return table.VerificationCodes
}

func (r VerificationCodeRepository) Create(ctx *appcontext.AppContext, code domain.VerificationCode) error {
	mapper := mapping.VerificationCodeMapper{}
	doc, err := mapper.FromDomainToModel(code)
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

func (r VerificationCodeRepository) Find(ctx *appcontext.AppContext, ip, email, code string) (*domain.VerificationCode, error) {
	var c = r.getTable()

	stmt := postgres.SELECT(
		c.AllColumns,
	).
		FROM(c).
		WHERE(
			c.IP.EQ(postgres.String(ip)).
				AND(c.Email.EQ(postgres.String(email))).
				AND(c.Code.EQ(postgres.String(code))),
		)

	var doc model.VerificationCodes
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.VerificationCodeMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r VerificationCodeRepository) TotalSentTodayByIp(ctx *appcontext.AppContext, ip string) (int64, error) {
	var (
		c            = r.getTable()
		startOfToday = manipulation.StartOfDay(manipulation.NowUTC())
	)

	stmt := postgres.SELECT(
		postgres.COUNT(c.ID).AS("count_result.total"),
	).
		FROM(c).
		WHERE(
			c.IP.EQ(postgres.String(ip)).
				AND(c.CreatedAt.GT_EQ(postgres.TimestampzT(startOfToday))),
		)

	var result = database.CountResult{}
	err := stmt.QueryContext(ctx.Context(), r.getDB(), &result)
	return result.Total, err
}

func (r VerificationCodeRepository) DeleteExpired(ctx *appcontext.AppContext) error {
	var (
		c   = r.getTable()
		now = manipulation.NowUTC()
	)

	stmt := c.DELETE().WHERE(c.ExpiresAt.LT(postgres.TimestampzT(now)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}
