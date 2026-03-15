package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormMemoRepository struct {
	db *gorm.DB
}

func NewGormMemoRepository(db *gorm.DB) *GormMemoRepository {
	return &GormMemoRepository{db: db}
}

func (r *GormMemoRepository) FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*entity.Memo, error) {
	tx := GetTx(ctx, r.db)
	var models []MemoModel
	if err := tx.Where("book_id = ?", bookID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Memo, len(models))
	for i := range models {
		entities[i] = toMemoEntity(&models[i])
	}
	return entities, nil
}

func (r *GormMemoRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Memo, error) {
	tx := GetTx(ctx, r.db)
	var model MemoModel
	if err := tx.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("memo not found")
		}
		return nil, err
	}
	return toMemoEntity(&model), nil
}

func (r *GormMemoRepository) Save(ctx context.Context, memo *entity.Memo) error {
	tx := GetTx(ctx, r.db)
	model := toMemoModel(memo)
	return tx.Save(model).Error
}

func (r *GormMemoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := GetTx(ctx, r.db)
	return tx.Delete(&MemoModel{}, "id = ?", id).Error
}
