package main

import "gorm.io/gorm"

type Enrollment struct {
	ID        uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"foreignKey:ID"`
	StudentID uint `gorm:"foreignKey:ID"`
	Status    string
	Grade     string
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}

func AddNewEnrollment(db *gorm.DB, enrollment *Enrollment) error {
	return db.Create(enrollment).Error
}

func RetrieveEnrollmentByID(db *gorm.DB, id uint) (*Enrollment, error) {
	var enrollment Enrollment
	if err := db.First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func UpdateEnrollment(db *gorm.DB, enrollment *Enrollment) error {
	return db.Save(enrollment).Error
}

func DeleteEnrollmentByID(db *gorm.DB, id uint) error {
	return db.Delete(&Enrollment{}, id).Error
}
