package db

import (
	"github.com/google/pprof/profile"
	"gorm.io/gorm"
)

type Grades struct {
	avg    uint
	grades [25]uint
}

type Student struct {
	gorm.Model
	Code       uint
	FullName   string
	Score      Grades
	Gmail      string
	ProfilePic string
}
type Teacher struct {
	Code       uint
	Gmail      string
	ProfilePic string
	FullName   string
}
