package repository

import (
	"context"
	"errors"
	"golang-ai/internal/entity"
	"golang-ai/internal/pkg/serverutils"
	"golang-ai/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type INotebookRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) INotebookRepository
	Create(ctx context.Context, notebook *entity.Notebook) error
	GetById(ctx context.Context, id uuid.UUID) (*entity.Notebook, error)
	Update(ctx context.Context, notebook *entity.Notebook) error
}

type notebookRepository struct {
	db database.DatabaseQueryer
}

func (n *notebookRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) INotebookRepository {
	return &notebookRepository{
		db: tx,
	}
}

func (n *notebookRepository) Create(ctx context.Context, notebook *entity.Notebook) error {
	_, err := n.db.Exec(
		ctx,
		`INSERT INTO notebook (id, name, parent_id, created_at, updated_at, deleted_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		notebook.Id,
		notebook.Name,
		notebook.ParentId,
		notebook.CreatedAt,
		notebook.UpdatedAt,
		notebook.DeletedAt,
		notebook.IsDeleted,
	)
	if err != nil {
		return err
	}

	return nil
}
func (n *notebookRepository) Update(ctx context.Context, notebook *entity.Notebook) error {
	_, err := n.db.Exec(
		ctx,
		`UPDATE notebook SET name = $2, parent_id = $3, updated_at = $4 WHERE id = $1`,
		notebook.Id,
		notebook.Name,
		notebook.ParentId,
		notebook.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
func (n *notebookRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Notebook, error) {
	var notebook entity.Notebook
	err := n.db.QueryRow(
		ctx,
		`SELECT id, name, parent_id, created_at, updated_at, deleted_at, is_deleted FROM notebook n WHERE n.is_deleted = false AND n.id = $1`,
		id,
	).Scan(
		&notebook.Id,
		&notebook.Name,
		&notebook.ParentId,
		&notebook.CreatedAt,
		&notebook.UpdatedAt,
		&notebook.DeletedAt,
		&notebook.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serverutils.ErrNotFound
		}
		return nil, err
	}

	return &notebook, nil
}
func NewNotebookRepository(db *pgxpool.Pool) INotebookRepository {
	return &notebookRepository{
		db: db,
	}
}
