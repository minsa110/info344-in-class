package tasks

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct { // defining struct for object-oriented programming
	Session        *mgo.Session // connection to the mongo db
	DatabaseName   string
	CollectionName string
}

func (ms *MongoStore) Insert(newtask *NewTask) (*Task, error) { // right of "Insert" is the same as interface
	// invoking the method in the structure: MongoStore
	t := newtask.ToTask()
	t.ID = bson.NewObjectId()
	// assuming that active session has been started
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(t)
	// ^ can pass multiple parameters e.g. Insert(t1, t2)
	return t, err // the "caller"" will handle the error
}

func (ms *MongoStore) Get(ID interface{}) (*Task, error) {
	task := &Task{} // empty Task struct, so that mongo knows what to return
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task)
	// "One()" to get one thing out of the db
	return task, err
}
