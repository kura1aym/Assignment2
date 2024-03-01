package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres dbname=postgres password=postgres port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Department{}, &Course{}, &Instructor{}, &Student{}, &Enrollment{})
	if err != nil {
		return err
	}

	err = db.Exec("ALTER TABLE students ADD COLUMN age INT").Error
	if err != nil {
		return err
	}

	return nil
}

func GetRowCount(db *gorm.DB, tableName string) (int64, error) {
	var count int64
	if err := db.Model(&gorm.Model{}).Table(tableName).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetStudentsByDepartment(db *gorm.DB, departmentName string) ([]Student, error) {
	var students []Student
	if err := db.Table("students").Joins("JOIN enrollments ON students.id = enrollments.student_id").
		Joins("JOIN courses ON enrollments.course_id = courses.id").
		Joins("JOIN departments ON courses.department_id = departments.id").
		Where("departments.name = ?", departmentName).
		Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func GetCoursesByInstructor(db *gorm.DB, instructorName string) ([]Course, error) {
	var courses []Course
	if err := db.Table("courses").Joins("JOIN instructors ON courses.instructor_id = instructors.id").
		Where("instructors.name = ?", instructorName).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func GetEnrollmentsByStudent(db *gorm.DB, studentID uint) ([]Enrollment, error) {
	var enrollments []Enrollment
	if err := db.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func EnrollStudentInCourse(db *gorm.DB, studentID, courseID uint) error {
	tx := db.Begin()

	enrollment := Enrollment{StudentID: studentID, CourseID: courseID, Status: "Enrolled"}
	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to enroll student in course: %w", err)
	}

	tx.Commit()

	return nil
}

func GetStudentCountByDepartment(db *gorm.DB, departmentID uint) (int64, error) {
	var count int64
	if err := db.Model(&Student{}).
		Where("department_id = ?", departmentID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetCoursesWithEnrollmentCount(db *gorm.DB) (map[string]int, error) {
	var courses []struct {
		Title           string
		EnrollmentCount int
	}
	if err := db.Model(&Course{}).
		Select("courses.title, COUNT(enrollments.id) as enrollment_count").
		Joins("LEFT JOIN enrollments ON courses.id = enrollments.course_id").
		Group("courses.id, courses.title").
		Find(&courses).Error; err != nil {
		return nil, err
	}

	courseEnrollmentCounts := make(map[string]int)
	for _, course := range courses {
		courseEnrollmentCounts[course.Title] = course.EnrollmentCount
	}
	return courseEnrollmentCounts, nil
}

// go run .
func main() {
	db, err := ConnectDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	err = MigrateDB(db)
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	//CREATE
	student := &Student{Name: "Nurkhat", Surname: "Muratkhan", GPA: 3.4}
	if err := AddNewStudent(db, student); err != nil {
		panic(fmt.Errorf("failed to add student: %v", err))
	}
	fmt.Println("New student added successfully")

	//READ
	courseID := uint(1)
	course, err := RetrieveCourseByID(db, courseID)
	if err != nil {
		panic("failed to retrieve course: " + err.Error())
	}
	fmt.Println("Retrieved course:", course)

	//UPDATE
	instructorID := uint(3)
	instructor, err := RetrieveInstructorByID(db, instructorID)
	if err != nil {
		panic("failed to retrieve instructor: " + err.Error())
	}

	instructor.Salary = 600000
	err = UpdateInstructor(db, instructor)
	if err != nil {
		panic("failed to update instructor: " + err.Error())
	}
	fmt.Println("Instructor updated successfully")

	//DELETE
	departmentID := uint(3)
	err = DeleteDepartmentByID(db, departmentID)
	if err != nil {
		panic("failed to delete department: " + err.Error())
	}
	fmt.Println("Department deleted successfully")

	//QUERYING
	departmentsCount, err := GetRowCount(db, "departments")
	if err != nil {
		panic("failed to get row count for departments: " + err.Error())
	}

	coursesCount, err := GetRowCount(db, "courses")
	if err != nil {
		panic("failed to get row count for courses: " + err.Error())
	}

	instructorsCount, err := GetRowCount(db, "instructors")
	if err != nil {
		panic("failed to get row count for instructors: " + err.Error())
	}

	studentsCount, err := GetRowCount(db, "students")
	if err != nil {
		panic("failed to get row count for students: " + err.Error())
	}

	enrollmentsCount, err := GetRowCount(db, "enrollments")
	if err != nil {
		panic("failed to get row count for enrollments: " + err.Error())
	}

	fmt.Printf("Number of rows in departments table: %d\n", departmentsCount)
	fmt.Printf("Number of rows in courses table: %d\n", coursesCount)
	fmt.Printf("Number of rows in instructors table: %d\n", instructorsCount)
	fmt.Printf("Number of rows in students table: %d\n", studentsCount)
	fmt.Printf("Number of rows in enrollments table: %d\n", enrollmentsCount)

	students, err := GetStudentsByDepartment(db, "Computer Science")
	if err != nil {
		fmt.Println("Failed to retrieve students:", err)
		return
	}

	fmt.Println("Students in Computer Science department:")
	for _, student := range students {
		fmt.Printf("ID: %d, Name: %s, Surname: %s, GPA: %.2f\n", student.ID, student.Name, student.Surname, student.GPA)
	}

	courseList, err := GetCoursesByInstructor(db, "John Doe")
	if err != nil {
		fmt.Println("Failed to retrieve courses:", err)
		return
	}

	fmt.Println("Courses taught by John Doe:")
	for _, course := range courseList {
		fmt.Printf("ID: %d, Code: %s, Title: %s\n", course.ID, course.Code, course.Title)
	}

	enrollments, err := GetEnrollmentsByStudent(db, 103)
	if err != nil {
		fmt.Println("Failed to retrieve enrollments:", err)
		return
	}

	fmt.Println("Enrollments for student with ID 123:")
	for _, enrollment := range enrollments {
		fmt.Printf("ID: %d, CourseID: %d, StudentID: %d, Status: %s, Grade: %s\n", enrollment.ID, enrollment.CourseID, enrollment.StudentID, enrollment.Status, enrollment.Grade)
	}
	
	//TRANSACTION
	err = EnrollStudentInCourse(db, 104, 2)
	if err != nil {
		panic("failed to enroll student in course: " + err.Error())
	}

	fmt.Println("Student enrolled in course successfully!")

	//HOOKS
	student2 := &Student{ID: 109, Name: "John", Surname: "Doe", GPA: 3.5}
	if err := db.Create(student2).Error; err != nil {
		panic(fmt.Errorf("failed to create student: %v", err))
	}
	fmt.Println("New student created successfully")

	instructor2 := &Instructor{ID: 1, Name: "Jane", Surname: "Smith", Salary: 50000}
	if err := db.Save(instructor2).Error; err != nil {
		panic(fmt.Errorf("failed to update instructor: %v", err))
	}
	fmt.Println("Instructor updated successfully")

	//SOFT DELETE
	studentID := uint(106)
	if err := SoftDeleteStudentByID(db, studentID); err != nil {
		panic(err)
	}

	fmt.Println("Student soft deleted successfully")

	//CUSTOM QUERIES
	departmentID2 := uint(1)
	studentCount, err := GetStudentCountByDepartment(db, departmentID2)
	if err != nil {
		panic("failed to get student count by department: " + err.Error())
	}
	department, err := RetrieveDepartmentByID(db, departmentID2)
	if err != nil {
		panic("failed to retrieve department: " + err.Error())
	}
	fmt.Printf("Total number of students in department '%s': %d\n", department.Name, studentCount)

	courses, err := GetCoursesWithEnrollmentCount(db)
	if err != nil {
		panic("failed to get courses with enrollment count: " + err.Error())
	}
	fmt.Println("Courses with enrollment count:")
	for title, enrollmentCount := range courses {
		fmt.Printf("Course: %s, Enrollment Count: %d\n", title, enrollmentCount)
	}
}

//student := Student{Name: "John", Surname: "Doe", GPA: 3.5}
//db.Create(&student)
//
//var retrievedStudent Student
//db.First(&retrievedStudent, "name = ?", "John")
//fmt.Println("Retrieved Student:", retrievedStudent)
//
//db.Model(&retrievedStudent).Update("GPA", 4.0)
//
//db.Delete(&retrievedStudent)

//result := testing.Main(func(pat, str string) (bool, error) { return true, nil },
//	// Передаем ваши тестовые функции
//	TestGetStudentCountByDepartment)
//
//// Если результат тестирования не успешен, выходим с кодом ошибки
//if !result {
//	os.Exit(1)
//}

//departments := []Department{
//	{Name: "Computer Science", Dean: "John Doe"},
//	{Name: "Mathematics", Dean: "Jane Smith"},
//	{Name: "Physics", Dean: "Michael Johnson"},
//}
//
//for i := range departments {
//	db.Create(&departments[i])
//}
//
//instructors := []Instructor{
//	{Name: "Aliya", Surname: "Maxatkyzy", Salary: 500000},
//	{Name: "Aidar", Surname: "Bakytzhanuly", Salary: 550000},
//	{Name: "Zhanna", Surname: "Sagintayeva", Salary: 520000},
//}
//
//for i := range instructors {
//	db.Create(&instructors[i])
//}
//
//students := []Student{
//	{ID: 101, Name: "Kuralay", Surname: "Mukhtar", GPA: 3.85},
//	{ID: 102, Name: "Bolat", Surname: "Maratuly", GPA: 2.84},
//	{ID: 103, Name: "Makpal", Surname: "Berikkyzy", GPA: 1.8},
//}
//
//for i := range students {
//	formattedGPA := students[i].FormattedGPA()
//	gpaFloat, err := strconv.ParseFloat(formattedGPA, 64)
//	if err != nil {
//		fmt.Println("Error converting GPA to float64:", err)
//		continue
//	}
//	students[i].GPA = gpaFloat
//
//	db.Create(&students[i])
//}
//
//courses := []Course{
//	{Title: "Introduction to Computer Science", Code: "CS101", Department: departments[0], Instructor: instructors[0]},
//	{Title: "Calculus I", Code: "MATH101", Department: departments[1], Instructor: instructors[1]},
//	{Title: "Classical Mechanics", Code: "PHY101", Department: departments[2], Instructor: instructors[2]},
//}
//
//for i := range courses {
//	db.Create(&courses[i])
//}
//
//enrollments := []Enrollment{
//	{CourseID: courses[0].ID, StudentID: students[0].ID, Status: "Approved", Grade: "A"},
//	{CourseID: courses[1].ID, StudentID: students[1].ID, Status: "Approved", Grade: "B"},
//	{CourseID: courses[2].ID, StudentID: students[2].ID, Status: "Approved", Grade: "C"},
//}
//
//for i := range enrollments {
//	db.Create(&enrollments[i])
//}
