package repository

import (
	"context"
	"fmt"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
)

type MovimentacaoRepository struct {
	db DBTX
}

func NewMovimentacaoRepository(db DBTX) *MovimentacaoRepository {
	return &MovimentacaoRepository{
		db: db,
	}
}

func (r *MovimentacaoRepository) Create(ctx context.Context, movimentacao *model.Movimentacao) error {
	query := `
	INSERT INTO movimentacoes (item_id, natureza, quantidade)
	VALUES ($1, $2, $3)
	RETURNING id, data
	`
	err := r.db.QueryRow(
		ctx,
		query,
		movimentacao.ItemID,
		movimentacao.Natureza,
		movimentacao.Quantidade,
	).Scan(
		&movimentacao.ID,
		&movimentacao.Data,
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar movimentacao: %w", err)
	}

	return nil
}

func (r *MovimentacaoRepository) FindByItemId(ctx context.Context, itemID int64) ([]model.Movimentacao, error) {
	query := `
	SELECT id, item_id, natureza, quantidade, data
	FROM movimentacoes
	WHERE item_id = $1
	ORDER BY data DESC
	`
	rows, err := r.db.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentacoes: %w", err)
	}
	defer rows.Close()

	var movimentacoes []model.Movimentacao

	for rows.Next() {
		var movimentacao model.Movimentacao
		err := rows.Scan(
			&movimentacao.ID,
			&movimentacao.ItemID,
			&movimentacao.Natureza,
			&movimentacao.Quantidade,
			&movimentacao.Data,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler movimentacao: %w", err)
		}
		movimentacoes = append(movimentacoes, movimentacao)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer movimentacoes: %w", err)
	}

	return movimentacoes, nil
}

func (r *MovimentacaoRepository) FindAll(ctx context.Context) ([]model.Movimentacao, error) {
	query := `
	SELECT id, item_id, natureza, quantidade, data
	FROM movimentacoes
	ORDER BY data DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentacoes: %w", err)
	}
	rows.Close()

	var movimentacoes []model.Movimentacao

	for rows.Next() {
		var movimentacao model.Movimentacao

		err := rows.Scan(
			&movimentacao.ID,
			&movimentacao.ItemID,
			&movimentacao.Natureza,
			&movimentacao.Quantidade,
			&movimentacao.Data,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler movimentacao: %w", err)
		}

		movimentacoes = append(movimentacoes, movimentacao)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer movimentacoes: %w", err)
	}

	return movimentacoes, nil
}
