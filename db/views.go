package db

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetStudent() Student {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect dateabase")
	}
	ctx := context.Background()
	student, err := gorm.G[Student](db).First(ctx)
	return student
}

func TestGetStudentGrades() Grades {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect dateabase")
	}
	ctx := context.Background()
	grades, err := gorm.G[Grades](db).Where("Student_Code = ?", 0).Find(ctx)
	return grades
}
