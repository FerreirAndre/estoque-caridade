package service

import (
	"context"
	"errors"
	"strings"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
	"github.com/FerreirAndre/estoque-caridade/internal/repository"
)

var (
	ErrItemNotFound        = errors.New("item nao encontrado")
	ErrInvalidName         = errors.New("nome do item invalido")
	ErrInvalidUnit         = errors.New("unidade de medida invalida")
	ErrInvalidMinimumStock = errors.New("estoque minimo nao pode ser negativo")
)

type ItemService struct {
	repository *repository.ItemRepository
}

func NewItemService(repository *repository.ItemRepository) *ItemService {
	return &ItemService{
		repository: repository,
	}
}

func (s *ItemService) Create(ctx context.Context, item *model.Item) error {
	item.Nome = strings.TrimSpace(item.Nome)

	if item.Nome == "" {
		return ErrInvalidName
	}

	if !isValidUnit(item.Unidade) {
		return ErrInvalidUnit
	}

	if item.EstoqueMinimo < 0 {
		return ErrInvalidMinimumStock
	}

	item.Quantidade = 0

	return s.repository.Create(ctx, item)
}

func (s *ItemService) FindById(ctx context.Context, id int64) (*model.Item, error) {
	item, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, ErrItemNotFound
	}

	return item, nil
}

func (s *ItemService) FindAll(ctx context.Context) ([]model.Item, error) {
	return s.repository.FindAll(ctx)
}

func (s *ItemService) Update(ctx context.Context, item *model.Item) error {
	item.Nome = strings.TrimSpace(item.Nome)

	if item.Nome == "" {
		return ErrInvalidName
	}

	if !isValidUnit(item.Unidade) {
		return ErrInvalidUnit
	}

	if item.EstoqueMinimo < 0 {
		return ErrInvalidMinimumStock
	}

	existing, err := s.FindById(ctx, item.ID)
	if err != nil {
		return err
	}

	item.Quantidade = existing.Quantidade

	return s.repository.Update(ctx, item)
}

func (s *ItemService) Delete(ctx context.Context, id int64) error {
	_, err := s.FindById(ctx, id)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, id)
}

func isValidUnit(unit model.UnidadeMedida) bool {
	switch unit {
	case model.KG, model.PACOTE, model.LITRO:
		return true
	default:
		return false
	}
}
