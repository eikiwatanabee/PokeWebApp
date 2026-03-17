package persistence

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTeamRepository struct {
	db *gorm.DB
}

func NewGormTeamRepository(db *gorm.DB) *GormTeamRepository {
	return &GormTeamRepository{db: db}
}

func (r *GormTeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	db := GetTx(ctx, r.db)
	var model TeamModel
	if err := db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return mapTeamModelToEntity(&model), nil
}

func (r *GormTeamRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.Team, error) {
	db := GetTx(ctx, r.db)
	var models []TeamModel
	if err := db.Where("tenant_id = ?", tenantID).Find(&models).Error; err != nil {
		return nil, err
	}
	teams := make([]*entity.Team, len(models))
	for i, m := range models {
		teams[i] = mapTeamModelToEntity(&m)
	}
	return teams, nil
}

func (r *GormTeamRepository) Save(ctx context.Context, team *entity.Team) error {
	db := GetTx(ctx, r.db)
	model := mapTeamEntityToModel(team)
	return db.Save(model).Error
}

func (r *GormTeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := GetTx(ctx, r.db)
	return db.Delete(&TeamModel{}, "id = ?", id).Error
}
