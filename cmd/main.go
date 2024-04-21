package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var users []User
var maxId uint64

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
	maxId = 3
}

func main() {
	http.HandleFunc("/users", UserServer)
	log.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetAllUsers(w)
	case http.MethodPost:
		decode := json.NewDecoder(r.Body) //Obtengo el body del post que me llega
		var u User
		if err := decode.Decode(&u); err != nil {
			MsgResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		fmt.Printf("%+v\n", u)
		PostUser(w, u)
	default:
		InvalidMethod(w)
	}
}

func GetAllUsers(w http.ResponseWriter) {
	DataResponse(w, http.StatusOK, users)
}

func PostUser(w http.ResponseWriter, data interface{}) {
	user := data.(User) //Casteo la interfaz a User

	if user.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "First name is required")
		return
	}

	if user.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "Last name is required")
		return
	}

	if user.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "Email is required")
		return
	}

	maxId++
	user.ID = maxId
	users = append(users, user)
	DataResponse(w, http.StatusCreated, user)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}

func MsgResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": %s}`, status, message)
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)

}
