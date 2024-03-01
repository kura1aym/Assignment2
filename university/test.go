package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestGetStudentCountByDepartment(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	db.AutoMigrate(&Department{})
	db.Create(&Department{ID: 1, Name: "Test Department"})

	db.AutoMigrate(&Student{})
	db.Create(&Student{ID: 9, Name: "Test Student", Surname: "Test", DepartmentID: 3})

	count, err := GetStudentCountByDepartment(db, 3)
	if err != nil {
		t.Fatalf("failed to get student count by department: %v", err)
	}

	expected := int64(3)
	if count != expected {
		t.Errorf("got %d students in 'Test Department', expected %d", count, expected)
	}
}
