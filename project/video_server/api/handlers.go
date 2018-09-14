package main 

import(
     "io"
     "net/http" 
     "github.com/julienschmidt/httprouter"
)

func CreateUser() {
	io.WriteString(w,"Create User Handler")
}