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
	err = gorm.G[Student](db).Create(ctx, &Student{Code: 0, FullName: "Victor Bivol", Gmail: "vbivol394@gmail.com", ProfilePic: "https://www.napoleon.org/wp-content/thumbnails/uploads/2006/01/t5_09ter_napoleon_fontainebleau_full-tt-width-637-height-911-crop-1-bgcolor-ffffff-lazyload-0.jpg", Avg: 9, Grades: []uint{10}})
	if err != nil {
		panic("failed to create student")
	}
	fmt.Println("totul a mers bine")
}
