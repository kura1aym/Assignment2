package main

import (
	"gorm.io/gorm"
	"time"
)

type Instructor struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Surname   string
	Salary    int
	Courses   []Course
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (Instructor) TableName() string {
	return "instructors"
}

func (i *Instructor) BeforeUpdate(tx *gorm.DB) (err error) {
	i.UpdatedAt = time.Now()
	return nil
}

func AddNewInstructor(db *gorm.DB, instructor *Instructor) error {
	return db.Create(instructor).Error
}

func RetrieveInstructorByID(db *gorm.DB, id uint) (*Instructor, error) {
	var instructor Instructor
	if err := db.First(&instructor, id).Error; err != nil {
		return nil, err
	}
	return &instructor, nil
}

func UpdateInstructor(db *gorm.DB, instructor *Instructor) error {
	return db.Save(instructor).Error
}

func DeleteInstructorByID(db *gorm.DB, id uint) error {
	return db.Delete(&Instructor{}, id).Error
}
