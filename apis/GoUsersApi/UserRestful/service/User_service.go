package service

import (
	"DataRestful/users"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"DataRestful/mongo"
)


type UserService struct {
	collection *mgo.Collection
	hash       users.Hash

}

func userModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewUserService(session *mongo.Session, dbName string, collectionName string, hash users.Hash) UserService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(userModelIndex())
	return UserService{collection, hash}
}

func (p *UserService) Create(u *users.User) error {
	//user := newUserModel(u)
	hashedPassword, err := p.hash.Generate(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return p.collection.Insert(&u)
}

func (p *UserService) GetByUsername(username string) (*users.User, error) {
	model := &users.User{}
	err := p.collection.Find(bson.M{"username": username}).One(&model)
	return model, err
}