package main

import "database/sql"
import "log"
import _ "github.com/go-sql-driver/mysql"

func Authenticate(user Cred) string {
	db, err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/test1")
    if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	results, err := db.Query("SELECT * FROM AUTHENTICATION where username='"+user.Username+"' and password= SHA1('"+user.Password+"')")
	if err != nil {
		panic(err.Error())
	}
	if(results.Next()){
		return "VALID"
	}
	return "INVALID"

}

func Register(user Cred) Cred{
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test1")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
    insert, err := db.Query("INSERT INTO AUTHENTICATION VALUES ('"+user.Username+"', SHA1('"+user.Password+"') )")
    if err != nil {
        panic(err.Error())
    }
    defer insert.Close()
    return user
}