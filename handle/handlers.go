package handle

import (
	"encoding/json"
	"example/hello/db"
	"net/http"
)

type StudentJson struct {
	Code       uint   `json:"Code"`
	FullName   string `json:"FullName"`
	Gmail      string `json:"Gamil"`
	ProfilePic string `json:"ProfilePic"`
	Avg        uint   `json: Avg`
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	student := db.TestGetStudent
	s := StudentJson{
		Code:       student().Code,
		FullName:   student().FullName,
		Gmail:      student().Gmail,
		ProfilePic: student().ProfilePic,
		Avg:        student().Avg,
	}
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "Failed to encode json", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
