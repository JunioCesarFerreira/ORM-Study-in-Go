package repository

import (
	"m/tests/ClassWithGorm/entities"

	"gorm.io/gorm"
)

func InsertClass(db *gorm.DB, name string) (*entities.Class, error) {
	class := entities.Class{Name: name}
	if err := db.Create(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func ReadClass(db *gorm.DB, classID int) (*entities.Class, error) {
	var class entities.Class
	err := db.Preload("Objects.Items").First(&class, classID).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func UpdateClass(db *gorm.DB, classID int, newName string) error {
	return db.Model(&entities.Class{}).Where("id = ?", classID).Update("name", newName).Error
}

func DeleteClass(db *gorm.DB, classID int) error {
	return db.Delete(&entities.Class{}, classID).Error
}
