package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "io"
)

func loginpage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the Login Page!!\n")
	fmt.Println("Endpoint hit : homepage")
}

func AccountVerify(w http.ResponseWriter, r *http.Request){   
    var user Cred
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) 
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    a:=Authenticate(user)
    if(a=="VALID"){
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusAccepted)
    }else{
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusUnauthorized)
    }

    if err := json.NewEncoder(w).Encode(a); err != nil {
        panic(err)
    }
}

func AccountCreate(w http.ResponseWriter, r *http.Request){
    var user Cred
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) 
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    a := Register(user)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(a); err != nil {
        panic(err)
    }
}