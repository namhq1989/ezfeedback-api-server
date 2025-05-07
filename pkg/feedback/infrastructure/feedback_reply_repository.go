package infrastructure

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/table"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/infrastructure/mapping"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackReplyRepository struct {
	db *database.Database
}

func NewFeedbackReplyRepository(db *database.Database) FeedbackReplyRepository {
	r := FeedbackReplyRepository{
		db: db,
	}

	return r
}

func (r FeedbackReplyRepository) getDB() *sql.DB {
	return r.db.GetPgDb()
}

func (FeedbackReplyRepository) getTable() *table.FeedbackRepliesTable {
	return table.FeedbackReplies
}

func (r FeedbackReplyRepository) Create(ctx *appcontext.AppContext, reply domain.FeedbackReply) error {
	mapper := mapping.FeedbackReplyMapper{}
	doc, err := mapper.FromDomainToModel(reply)
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

func (r FeedbackReplyRepository) Update(ctx *appcontext.AppContext, reply domain.FeedbackReply) error {
	mapper := mapping.FeedbackReplyMapper{}
	doc, err := mapper.FromDomainToModel(reply)
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

func (r FeedbackReplyRepository) Delete(ctx *appcontext.AppContext, reply domain.FeedbackReply) error {
	var (
		fr = r.getTable()
	)

	stmt := fr.DELETE().WHERE(fr.ID.EQ(postgres.String(reply.ID)))
	_, err := stmt.ExecContext(ctx.Context(), r.getDB())
	return err
}

func (r FeedbackReplyRepository) FindByID(ctx *appcontext.AppContext, replyID string) (*domain.FeedbackReply, error) {
	if !uuid.IsValidID(replyID) {
		return nil, apperrors.Feedback.InvalidReply
	}

	var fr = r.getTable()

	stmt := postgres.SELECT(
		fr.AllColumns,
	).
		FROM(fr).
		WHERE(fr.ID.EQ(postgres.String(replyID)))

	var doc model.FeedbackReplies
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &doc); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper    = mapping.FeedbackReplyMapper{}
		result, _ = mapper.FromModelToDomain(doc)
	)
	return result, nil
}

func (r FeedbackReplyRepository) FindWithFilter(ctx *appcontext.AppContext, filter domain.FeedbackReplyFilter) ([]domain.FeedbackReply, error) {
	var (
		fr        = r.getTable()
		offset    = filter.Limit * filter.Page
		whereStmt = fr.FeedbackID.EQ(postgres.String(filter.FeedbackID))
	)

	stmt := postgres.SELECT(
		fr.ID, fr.UserID, fr.Content, fr.IsEdited, fr.CreatedAt, fr.UpdatedAt,
	).
		FROM(fr).
		WHERE(whereStmt).
		LIMIT(filter.Limit).
		OFFSET(offset).
		ORDER_BY(fr.CreatedAt.DESC())

	var (
		docs   = make([]model.FeedbackReplies, 0)
		result = make([]domain.FeedbackReply, 0)
	)
	if err := stmt.QueryContext(ctx.Context(), r.getDB(), &docs); err != nil {
		if r.db.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		mapper = mapping.FeedbackReplyMapper{}
	)
	for _, doc := range docs {
		reply, err := mapper.FromModelToDomain(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *reply)
	}
	return result, nil
}
