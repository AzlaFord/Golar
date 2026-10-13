package handle

import (
	"example/hello/db"
	"io"
	"net/http"
)

func GetStudent(w http.ResponseWriter, r *http.Request) {
	student := db.TestGetStudent
	io.WriteString(w, student().FullName)
}
