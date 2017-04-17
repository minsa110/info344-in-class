package handlers

import (
	"encoding/json"
	"net/http"

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
	}
}

//HandleSpecificTask will handle requests for the /v1/tasks/some-task-id resource
func (ctx *Context) HandleSpecificTask(w http.ResponseWriter, r *http.Request) {
}
