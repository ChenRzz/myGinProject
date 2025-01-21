package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	Create(db *gorm.DB, entity *T) error
	Update(db *gorm.DB, entity *T) error
	Delete(db *gorm.DB, id uint) error
	FindByID(db *gorm.DB, id uint) (*T, error)
}

type BaseRepo[T any] struct {
}

func (b *BaseRepo[T]) Create(db *gorm.DB, entity *T) error {
	err := db.Create(entity).Error
	if err != nil {
		return errors.New("Failed to create entity")
	}
	return nil
}

func (b *BaseRepo[T]) FindByID(db *gorm.DB, id uint) (*T, error) {
	var entity T
	err := db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, errors.New("entity not found")
	}
	return &entity, nil
}

func (b *BaseRepo[T]) Update(db *gorm.DB, entity *T) error {
	err := db.Save(entity).Error
	if err != nil {
		return errors.New("Failed to update entity")
	}
	return nil
}

func (b *BaseRepo[T]) Delete(db *gorm.DB, id uint) error {
	var entity T
	err := db.Model(entity).Updates(
		map[string]interface{}{
			"is_deleted": gorm.Expr("id"),
			"deleted_at": time.Now(),
		}).Error
	if err != nil {
		return errors.New("Failed to update entity")
	}
	return nil
}
