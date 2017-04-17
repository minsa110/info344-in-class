package tasks

//Store defines an abstract interface for a Task object store
type Store interface {
	//Insert inserts a NewTask and
	//returns the fully-populated Task or an error
	Insert(newtask *NewTask) (*Task, error)
	Get(ID interface{}) (*Task, error)
}

// Interfaces
// 	* guaranteed set of methods if imported
// 	* duct-typing (no need to "implement" the interface)
//				  (see mongostore.go, has all the same methods)
