package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Surname      string
	GPA          float64
	Courses      []Course `gorm:"many2many:enrollments"`
	DepartmentID uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *gorm.DeletedAt `gorm:"index"`
}

func (Student) TableName() string {
	return "students"
}

func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func SoftDeleteStudentByID(db *gorm.DB, studentID uint) error {
	var student Student
	if err := db.First(&student, studentID).Error; err != nil {
		return fmt.Errorf("failed to find student: %v", err)
	}

	if err := db.Delete(&student).Error; err != nil {
		return fmt.Errorf("failed to soft delete student: %v", err)
	}

	return nil
}

func (s *Student) FormattedGPA() string {
	if s.GPA > 4.0 {
		s.GPA = 4.0
	} else if s.GPA < 0 {
		s.GPA = 0
	}
	return fmt.Sprintf("%.2f", s.GPA)
}

func AddNewStudent(db *gorm.DB, student *Student) error {
	return db.Create(student).Error
}

func RetrieveStudentByID(db *gorm.DB, id uint) (*Student, error) {
	var student Student
	if err := db.First(&student, id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func UpdateStudent(db *gorm.DB, student *Student) error {
	return db.Save(student).Error
}

func DeleteStudentByID(db *gorm.DB, id uint) error {
	return db.Delete(&Student{}, id).Error
}
