package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
	"github.com/FerreirAndre/estoque-caridade/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidQuantity   = errors.New("quantidade deve ser maior que zero")
	ErrInvalidNature     = errors.New("natureza de movimentação inválida")
	ErrInsufficientStock = errors.New("estoque insuficiente")
)

type MovimentacaoService struct {
	pool *pgxpool.Pool
}

func NewMovimentacaoService(pool *pgxpool.Pool) *MovimentacaoService {
	return &MovimentacaoService{
		pool: pool,
	}
}

func (s *MovimentacaoService) Create(ctx context.Context, movimentacao *model.Movimentacao) error {
	if movimentacao.Quantidade <= 0 {
		return ErrInvalidQuantity
	}

	if movimentacao.Natureza != model.NaturezaEntrada &&
		movimentacao.Natureza != model.NaturezaSaida {
		return ErrInvalidNature
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	itemRepo := repository.NewItemRepository(tx)
	movimentacaoRepo := repository.NewMovimentacaoRepository(tx)

	item, err := itemRepo.FindByIdForUpdate(ctx, movimentacao.ItemID)
	if err != nil {
		return err
	}

	if item == nil {
		return ErrItemNotFound
	}

	var novaQuantidade int64

	switch movimentacao.Natureza {
	case model.NaturezaEntrada:
		if item.Quantidade > math.MaxInt64-movimentacao.Quantidade {
			return errors.New("quantidade final excede o limite permitido")
		}

		novaQuantidade = item.Quantidade + movimentacao.Quantidade

	case model.NaturezaSaida:
		if movimentacao.Quantidade > item.Quantidade {
			return ErrInsufficientStock
		}

		novaQuantidade = item.Quantidade - movimentacao.Quantidade
	}

	err = itemRepo.UpdateQuantity(ctx, item.ID, novaQuantidade)
	if err != nil {
		return err
	}

	err = movimentacaoRepo.Create(ctx, movimentacao)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar movimentação: %w", err)
	}

	return nil
}

func (s *MovimentacaoService) FindByItemId(ctx context.Context, itemID int64) ([]model.Movimentacao, error) {
	return repository.NewMovimentacaoRepository(s.pool).
		FindByItemId(ctx, itemID)
}

func (s *MovimentacaoService) FindAll(ctx context.Context) ([]model.Movimentacao, error) {
	return repository.NewMovimentacaoRepository(s.pool).
		FindAll(ctx)
}
