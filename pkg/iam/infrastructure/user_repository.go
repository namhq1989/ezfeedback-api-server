package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	r := UserRepository{
		db: db,
	}

	return r
}

func (r UserRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (UserRepository) getTable() *table.UsersTable {
	return table.Users
}

func (r UserRepository) Create(ctx *appcontext.AppContext, user domain.User) error {
	mapper := mapping.UserMapper{}
	doc, err := mapper.FromDomainToModel(user)
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

func (r UserRepository) Update(ctx *appcontext.AppContext, user domain.User) error {
	mapper := mapping.UserMapper{}
	doc, err := mapper.FromDomainToModel(user)
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

func (r UserRepository) FindByID(ctx *appcontext.AppContext, userID string) (*domain.User, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	var u = r.getTable()

	stmt := postgres.SELECT(
		u.AllColumns,
	).
		FROM(u).
		WHERE(u.ID.EQ(postgres.String(userID)))

	var doc model.Users
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.UserMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r UserRepository) FindByEmail(ctx *appcontext.AppContext, email string) (*domain.User, error) {
	var u = r.getTable()

	stmt := postgres.SELECT(
		u.AllColumns,
	).
		FROM(u).
		WHERE(u.Email.EQ(postgres.String(email)))

	var doc model.Users
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.UserMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}
