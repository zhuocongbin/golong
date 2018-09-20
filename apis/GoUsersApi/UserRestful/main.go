package main

import (

	"log"
	"DataRestful/mongo"
	"DataRestful/server"
	"DataRestful/crypto"
	"DataRestful/service"
)

func main() {
	ms, err := mongo.NewSession("127.0.0.1:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer ms.Close()

	h := crypto.Hash{}
	u := service.NewUserService(ms.Copy(), "go_web_server", "user", &h)
	s := server.NewServer(u)

	s.Start()
}
