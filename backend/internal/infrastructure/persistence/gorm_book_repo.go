package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) *GormBookRepository {
	return &GormBookRepository{db: db}
}

func (r *GormBookRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	tx := GetTx(ctx, r.db)
	var model BookModel
	if err := tx.Preload("Tags").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return toBookEntity(&model), nil
}

func (r *GormBookRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	tx := GetTx(ctx, r.db)
	var model BookModel
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Tags").
		First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return toBookEntity(&model), nil
}

func (r *GormBookRepository) FindByUserID(ctx context.Context, userID uuid.UUID, status string, tagID string, page, limit int) ([]*entity.Book, int64, error) {
	tx := GetTx(ctx, r.db)
	query := tx.Model(&BookModel{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if tagID != "" {
		query = query.Joins("JOIN book_tags ON book_tags.book_model_id = books.id").
			Where("book_tags.tag_model_id = ?", tagID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []BookModel
	offset := (page - 1) * limit
	if err := query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.Book, len(models))
	for i := range models {
		entities[i] = toBookEntity(&models[i])
	}
	return entities, total, nil
}

func (r *GormBookRepository) Save(ctx context.Context, book *entity.Book) error {
	tx := GetTx(ctx, r.db)
	model := toBookModel(book)

	// Upsert with tag association
	if err := tx.Save(model).Error; err != nil {
		return err
	}

	// Replace tag associations
	if err := tx.Model(model).Association("Tags").Replace(model.Tags); err != nil {
		return err
	}

	return nil
}

func (r *GormBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := GetTx(ctx, r.db)
	// Clear tag associations first
	model := &BookModel{ID: id}
	if err := tx.Model(model).Association("Tags").Clear(); err != nil {
		return err
	}
	return tx.Delete(model).Error
}
