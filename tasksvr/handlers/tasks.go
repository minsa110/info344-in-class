package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/minsa110/info344-in-class/tasksvr/models/tasks"
)

//HandleTasks will handle requests for the /v1/tasks resource
func (ctx *Context) HandleTasks(w http.ResponseWriter, r *http.Request) {
	// gets called upon "POST"
	switch r.Method {
	case "POST": // to add something to the collection
		decoder := json.NewDecoder(r.Body)
		newtask := &tasks.NewTask{}
		if err := decoder.Decode(newtask); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest) // http error: report to user
			return                                               // so that it doesn't keep going
		}

		// validate
		if err := newtask.Validate(); err != nil {
			http.Error(w, "error validating task: "+err.Error(), http.StatusBadRequest)
			return
		}

		task, err := ctx.TasksStore.Insert(newtask)
		if err != nil {
			http.Error(w, "error inserting task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// to encode back to the client
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(task)
	case "GET":
		tasks, err := ctx.TasksStore.GetAll()
		if err != nil {
			http.Error(w, "error getting tasks: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// this is same as above. make a method later!! ***
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(tasks)
	}
}

//HandleSpecificTask will handle requests for the /v1/tasks/some-task-id resource
func (ctx *Context) HandleSpecificTask(w http.ResponseWriter, r *http.Request) {
	_, id := path.Split(r.URL.Path) // split the last part of the URL

	switch r.Method {
	case "GET":
		task, err := ctx.TasksStore.Get(id)
		if err != nil {
			http.Error(w, "error finding task: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add(headerContentType, contentTypeJSONUTF8) // let reader know the response is in JSON
		encoder := json.NewEncoder(w)
		encoder.Encode(task)
	case "PATCH":
		decoder := json.NewDecoder(r.Body)
		task := &tasks.Task{}
		if err := decoder.Decode(task); err != nil {
			http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		task.ID = id // so that the user doesn't have to provide the id to update

		if err := ctx.TasksStore.Update(task); err != nil {
			http.Error(w, "error updating: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("update successful!"))
	}
}
