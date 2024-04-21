package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"las_name"`
	Email     string `json:"email"`
}

var users []User

func init() { //La funcion init nos sirve para inicializar valores, si nosotros en un package tenemos la funcion init, esto es lo primero que se ejecuta.
	users = []User{
		{
			ID:        1,
			FirstName: "German",
			LastName:  "Sanz",
			Email:     "germansanz@gmail.com",
		},
		{
			ID:        2,
			FirstName: "Jhon",
			LastName:  "Doe",
			Email:     "johndoe@gmail.com",
		},
		{
			ID:        3,
			FirstName: "Mister",
			LastName:  "Jagger",
			Email:     "jaggermister@gmail.com",
		},
	}
}

func main() {
	http.HandleFunc("/users", UserServer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserServer(w http.ResponseWriter, r *http.Request) {
	var status int
	switch r.Method {
	case http.MethodGet:
		GetAllUsers(w)
	case http.MethodPost:
		status = 200
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, "success in post")
	default:
		status = 404
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, "not found")
	}
}

func GetAllUsers(w http.ResponseWriter) {
	DataResponse(w, http.StatusOK, users)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, _ := json.Marshal(users)
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}
