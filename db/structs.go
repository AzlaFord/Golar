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
	Code       uint   `gorm:"primaryKey" json:"Code"`
	FullName   string `json:""`
	Gmail      string
	ProfilePic string
	Avg        uint
	Grades     []uint `gorm:"serializer:json"`
}
type Teacher struct {
	Code       uint `gorm:"primaryKey"`
	Gmail      string
	ProfilePic string
	FullName   string
}
