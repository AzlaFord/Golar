package db

import (
	"gorm.io/gorm"
)

type Class struct {
	Code     uint `gorm:"primaryKey"`
	Students []Student
	Teachers [1]Teacher
}

type Student struct {
	gorm.Model
	Code       uint `gorm:"primaryKey"`
	FullName   string
	Gmail      string
	ProfilePic string
	avg        uint
	grades     [25]uint
}
type Teacher struct {
	Code       uint `gorm:"primaryKey"`
	Gmail      string
	ProfilePic string
	FullName   string
}
