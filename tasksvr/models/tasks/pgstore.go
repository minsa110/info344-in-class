// store specifically for postgres (SQL db)
package tasks

import "database/sql"

type PGStore struct {
	DB *sql.DB
}

// A LOT OF CODE. Only good if many to many relationships (compared to mongo),
// since mongo doesn't have "join"s
func (ps *PGStore) Insert(newtask *NewTask) (*Task, error) {
	t := newtask.ToTask()    // convert to full task object
	tx, err := ps.DB.Begin() // beginning of transaction
	if err != nil {
		return nil, err
	}
	// to insert the task into database
	sql := `insert into tasks
    (title, createdAt, modifiedAt, complete)
    values ($1,$2,$3,$4)
    returning id` // do not need a semicolon, will be automatically added if not there
	row := tx.QueryRow(sql, t.Title, t.CreatedAt, t.ModifiedAt, t.Complete) // returns ONE row of results
	err = row.Scan(&t.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// to insert all tags in the database
	sql = `insert into tags(taskID, tag)
    values ($1,$2)`
	for _, tag := range t.Tags {
		_, err := tx.Exec(sql, t.ID, tag)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return t, nil
}

func (ps *PGStore) Get(ID interface{}) (*Task, error) {
	return nil, nil
}

func (ps *PGStore) GetAll() ([]*Task, error) {
	return nil, nil
}

func (ps *PGStore) Update(task *Task) error {
	return nil
}
