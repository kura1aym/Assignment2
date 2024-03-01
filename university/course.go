package main

import (
	"gorm.io/gorm"
)

type Course struct {
	ID           uint `gorm:"primaryKey"`
	Code         string
	Title        string
	Department   Department `gorm:"embedded;embeddedPrefix:department_"`
	InstructorID uint
	DeletedAt    *gorm.DeletedAt `gorm:"index"`
}

func (Course) TableName() string {
	return "courses"
}

func AddNewCourse(db *gorm.DB, course *Course) error {
	return db.Create(course).Error
}

func RetrieveCourseByID(db *gorm.DB, id uint) (*Course, error) {
	var course Course
	if err := db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func UpdateCourse(db *gorm.DB, course *Course) error {
	return db.Save(course).Error
}

func DeleteCourseByID(db *gorm.DB, id uint) error {
	return db.Delete(&Course{}, id).Error
}
