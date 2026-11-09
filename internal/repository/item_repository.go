package repository

import (
	"context"
	"fmt"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
	"github.com/jackc/pgx/v5"
)

type ItemRepository struct {
	db DBTX
}

func NewItemRepository(db DBTX) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) Create(ctx context.Context, item *model.Item) error {
	query := `
	INSERT INTO itens (nome, quantidade, unidade_medida)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRow(ctx, query, item.Nome, item.Quantidade, item.Unidade)
	if err != nil {
		return fmt.Errorf("Erro ao cadastrar item: %w", err)
	}

	return nil
}

func (r *ItemRepository) FindById(ctx context.Context, id int64) (*model.Item, error) {
	query := `
	SELECT id, nome, quantidade, unidade_medida
	FROM itens
	WHERE id = $1
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Nome,
		&item.Quantidade,
		&item.Unidade,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("erro ao buscar item: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) FindAll(ctx context.Context) ([]model.Item, error) {
	query := `SELECT id, nome, quantidade, unidade_medida
	FROM itens
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao lista itens: %w", err)
	}
	defer rows.Close()

	var itens []model.Item

	for rows.Next() {
		var item model.Item

		err := rows.Scan(
			&item.ID,
			&item.Nome,
			&item.Quantidade,
			&item.Unidade,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler item: %w", err)
		}

		itens = append(itens, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer itens: %w", err)
	}

	return itens, nil
}
