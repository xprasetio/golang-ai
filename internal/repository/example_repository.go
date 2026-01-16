package repository

import (
	"context"
	"golang-ai/internal/entity"
	"golang-ai/pkg/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IExampleRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) IExampleRepository
	Ping(ctx context.Context) (*entity.Example, error)
}

type exampleRepository struct {
	db database.DatabaseQueryer
}

func (n *exampleRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) IExampleRepository {
	return &exampleRepository{
		db: tx,
	}
}

func (n *exampleRepository) Ping(ctx context.Context) (*entity.Example, error) {
	row := n.db.QueryRow(
		ctx,
		`SELECT 'hello' AS "message"`,
	)
	var example entity.Example
	err := row.Scan(
		&example.Message,
	)
	if err != nil {
		return nil, err
	}

	return &example, nil
}

func NewExampleRepository(db *pgxpool.Pool) IExampleRepository {
	return &exampleRepository{
		db: db,
	}
}
