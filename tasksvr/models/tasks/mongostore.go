package tasks

import (
	"time"

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
	// type assertion
	if sID, ok := ID.(string); ok {
		ID = bson.ObjectIdHex(sID) // binary structure that identifies the object
	}
	task := &Task{} // empty Task struct, so that mongo knows what to return
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task)
	// "One()" to get one thing out of the db
	return task, err
}

func (ms *MongoStore) GetAll() ([]*Task, error) { // MongoStore as the "receiver"
	tasks := []*Task{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&tasks)
	// Find(nil) for no filtering; if want to query, then that's where we'd query to filter the data
	// place what is returned in "tasks"
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ms *MongoStore) Update(task *Task) error {
	if sID, ok := task.ID.(string); ok {
		task.ID = bson.ObjectIdHex(sID)
	}
	task.ModifiedAt = time.Now()
	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName)
	updates := bson.M{"$set": bson.M{"complete": task.Complete, "modifiedat": task.ModifiedAt}}
	// all property names should be lowercase!!
	return col.UpdateId(task.ID, updates)
}
