package db

import (
	"gorm.io/gorm"
)

type Class struct {
	Code     uint    `gorm:"primaryKey"`
	Student  Student `gorm:"references:Code"`
	Teachers Teacher
}

type Grades struct {
	gorm.Model
	Code         uint `gorm:"primaryKey"`
	Grade        uint
	Student_Code Student `gorm:"foreignKey:Code"`
}

type Student struct {
	gorm.Model
	Code       uint `gorm:"primaryKey"`
	FullName   string
	Gmail      string
	ProfilePic string
	Avg        uint
}
type Teacher struct {
	Code       uint `gorm:"primaryKey"`
	Gmail      string
	ProfilePic string
	FullName   string
}
