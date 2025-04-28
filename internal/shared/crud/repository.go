package crud

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) FindAll() ([]T, error) {
	var models []T
	err := r.db.Find(&models).Error
	return models, err
}

func (r *Repository[T]) FindById(id string) (T, error) {
	var model T
	err := r.db.First(&model, "id = ?", id).Error
	return model, err
}

func (r *Repository[T]) FindFirstBy(field string, value string) (*T, error) {
	var user T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository[T]) Preload(entity *T, fields []string, key string, value string) {
	db := r.db
	for _, field := range fields {
		db = db.Preload(field)
	}
	_ = db.First(entity, key+" = ?", value).Error
}

func (r *Repository[T]) Save(entity *T) (*T, error) {
	result := r.db.Save(&entity)
	return entity, result.Error
}

func (r *Repository[T]) Delete(entity *T, hard bool) error {
	var db = r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(entity).Error
}

func (r *Repository[T]) DeleteById(id string) error {
	var model T
	return r.db.Delete(&model, "id = ?", id).Error
}
