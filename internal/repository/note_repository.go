package repository

import (
	"context"
	"errors"
	"fmt"
	"golang-ai/internal/entity"
	"golang-ai/internal/pkg/serverutils"
	"golang-ai/pkg/database"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type INoteRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) INoteRepository
	Create(ctx context.Context, note *entity.Note) error
	GetById(ctx context.Context, id uuid.UUID) (*entity.Note, error)
	GetByIdNotebookIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Note, error)
	Update(ctx context.Context, note *entity.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByNotebookId(ctx context.Context, notebookId uuid.UUID) error
}

type noteRepository struct {
	db database.DatabaseQueryer
}

func (n *noteRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) INoteRepository {
	return &noteRepository{
		db: tx,
	}
}

func (n *noteRepository) Create(ctx context.Context, note *entity.Note) error {
	_, err := n.db.Exec(
		ctx,
		`INSERT INTO note (id, title, content, notebook_id, created_at, updated_at, deleted_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		note.Id,
		note.Title,
		note.Content,
		note.NotebookId,
		note.CreatedAt,
		note.UpdatedAt,
		note.DeletedAt,
		note.IsDeleted,
	)
	if err != nil {
		return err
	}

	return nil
}
func (n *noteRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Note, error) {
	{
		note := &entity.Note{}
		err := n.db.QueryRow(
			ctx,
			`SELECT id, title, content, notebook_id, created_at, updated_at, deleted_at, is_deleted FROM note WHERE id = $1 AND is_deleted = false`,
			id,
		).Scan(
			&note.Id,
			&note.Title,
			&note.Content,
			&note.NotebookId,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.DeletedAt,
			&note.IsDeleted,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, serverutils.ErrNotFound
			}
			return nil, err
		}

		return note, nil
	}
}
func (n *noteRepository) Update(ctx context.Context, note *entity.Note) error {
	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET title = $1, content = $2, notebook_id = $3, updated_at = $4 WHERE id = $5 AND is_deleted = false`,
		note.Title,
		note.Content,
		note.NotebookId,
		note.UpdatedAt,
		note.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
func (n *noteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET is_deleted = true, deleted_at = NOW() WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
func (n *noteRepository) DeleteByNotebookId(ctx context.Context, notebookId uuid.UUID) error {
	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET is_deleted = true, deleted_at = NOW() WHERE notebook_id = $1`,
		notebookId,
	)
	if err != nil {
		return err
	}
	return nil
}
func (n *noteRepository) GetByIdNotebookIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Note, error) {
	if len(ids) == 0 {
		return []*entity.Note{}, nil
	}
	idStr := make([]string, 0)
	for _, id := range ids {
		idStr = append(idStr, fmt.Sprintf("'%s'", id.String()))
	}
	idSqlFormat := strings.Join(idStr, ", ")
	rows, err := n.db.Query(
		ctx,
		fmt.Sprintf(`SELECT id, title, content, notebook_id, created_at, updated_at FROM note WHERE notebook_id IN (%s) AND is_deleted = false`, idSqlFormat),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*entity.Note
	for rows.Next() {
		note := &entity.Note{}
		err := rows.Scan(
			&note.Id,
			&note.Title,
			&note.Content,
			&note.NotebookId,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return notes, nil
}
func NewNoteRepository(db *pgxpool.Pool) INoteRepository {
	return &noteRepository{
		db: db,
	}
}
