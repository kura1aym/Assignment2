package main

import "gorm.io/gorm"

type Department struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Dean      string
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (Department) TableName() string {
	return "departments"
}

func AddNewDepartment(db *gorm.DB, department *Department) error {
	return db.Create(department).Error
}

func RetrieveDepartmentByID(db *gorm.DB, id uint) (*Department, error) {
	var department Department
	if err := db.First(&department, id).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func UpdateDepartment(db *gorm.DB, department *Department) error {
	return db.Save(department).Error
}

func DeleteDepartmentByID(db *gorm.DB, id uint) error {
	return db.Delete(&Department{}, id).Error
}
