package service

import (
	"log"
	"testing"
	"DataRestful/users"
	"DataRestful/util"
	"DataRestful/mongo"
)


const (
	mongoUrl = "localhost:27017"
	dbName = "test_db"
	userCollectionName = "user"
)

func Test_UserService(t *testing.T) {
	t.Run("CreateUser", createUser_should_insert_user_into_mongo)
}

func createUser_should_insert_user_into_mongo(t *testing.T) {
	//Arrange
	session, err := mongo.NewSession(mongoUrl)
	if(err != nil) {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	//defer session.Close()

	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	s:=session.Copy()
	mockHash := util.Hash{}
	userService := NewUserService(s, dbName, userCollectionName,&mockHash)

	//testUserId	:= "integration_test_id"
	testUsername := "integration_test_user2"
	testPassword := "integration_test_password"
	user := users.User{

		Username: testUsername,
		Password: testPassword }

	//Act
	err = userService.Create(&user)

	//Assert
	if(err != nil) {
		t.Error("Unable to create user: %s", err)
	}
	var results []users.User
	session.GetCollection(dbName,userCollectionName).Find(nil).All(&results)

	count := len(results)
	if(count != 1) {
		t.Error("Incorrect number of results. Expected `1`, got: `%i`", count)
	}
	if(results[0].Username != user.Username) {
		t.Error("Incorrect Username. Expected `%s`, Got: `%s`", testUsername, results[0].Username)
	}
}
