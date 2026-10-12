package db

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDb() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect dateabase")
	}
	ctx := context.Background()

	db.AutoMigrate(&Student{})
	err = gorm.G[Student](db).Create(ctx, &Student{Code: 0, Score: Grades{avg: 9, grades: [25]uint{1, 10, 10, 10, 10, 10}}})
	if err != nil {
		panic("failed to create student")
	}
	fmt.Println("totul a mers bine")
}
