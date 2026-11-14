package repository

import (
	"context"
	"errors"
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
	INSERT INTO itens (nome, quantidade, unidade_medida, estoque_minimo)
	VALUES ($1, $2, $3,$4)
	RETURNING id
	`

	err := r.db.QueryRow(ctx, query, item.Nome, item.Quantidade, item.Unidade, item.EstoqueMinimo).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("Erro ao cadastrar item: %w", err)
	}

	return nil
}

func (r *ItemRepository) FindById(ctx context.Context, id int64) (*model.Item, error) {
	query := `
	SELECT id, nome, quantidade, unidade_medida, estoque_minimo
	FROM itens
	WHERE id = $1
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Nome,
		&item.Quantidade,
		&item.Unidade,
		&item.EstoqueMinimo,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("erro ao buscar item: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) FindAll(ctx context.Context) ([]model.Item, error) {
	query := `SELECT id, nome, quantidade, unidade_medida, estoque_minimo
	FROM itens
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao lista itens: %w", err)
	}
	defer rows.Close()

	itens := make([]model.Item, 0)

	for rows.Next() {
		var item model.Item

		err := rows.Scan(
			&item.ID,
			&item.Nome,
			&item.Quantidade,
			&item.Unidade,
			&item.EstoqueMinimo,
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

func (r *ItemRepository) Update(ctx context.Context, item *model.Item) error {
	query := `
	UPDATE itens
	SET nome = $1,
		unidade_medida = $2,
		estoque_minimo = $3
	WHERE id = $4
	`

	result, err := r.db.Exec(ctx, query, item.Nome, item.Unidade, item.EstoqueMinimo, item.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar o item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := `
	DELETE FROM itens WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir o item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ItemRepository) FindByIdForUpdate(ctx context.Context, id int64) (*model.Item, error) {
	query := `
        SELECT id, nome, quantidade, unidade_medida, estoque_minimo
        FROM itens
        WHERE id = $1
        FOR UPDATE
    `

	var item model.Item

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Nome,
		&item.Quantidade,
		&item.Unidade,
		&item.EstoqueMinimo,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("erro ao buscar item: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) UpdateQuantity(ctx context.Context, id int64, quantity int64) error {
	query := `
        UPDATE itens
        SET quantidade = $1
        WHERE id = $2
    `

	result, err := r.db.Exec(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estoque: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
