package api

import (
	"TaskScheduler/internal/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

var (
	empty         = struct{}{}
	limit         = 50
	errNotValidID = fmt.Errorf("not valid id")
	errTaskSearch = fmt.Errorf("Задача не найдена")
)

func task(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if !idIsValid(id) {
		writeError(w, errNotValidID, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if task == nil {
		writeError(w, errTaskSearch, http.StatusBadRequest)
		return
	}
	writeJson(w, task)

}
func addTask(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var buf bytes.Buffer
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if len(task.Title) == 0 {
		writeError(w, fmt.Errorf("Task is required"), http.StatusBadRequest)
		return
	}
	if err = checkDate(task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(task)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	response := struct {
		Id int `json:"id"`
	}{Id: int(id)}
	writeJson(w, response)
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var buf bytes.Buffer
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if len(task.Title) == 0 {
		writeError(w, fmt.Errorf("Title is required"), http.StatusBadRequest)
		return
	}
	if err := checkDate(task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(task); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJson(w, empty)
}

// The function tasksHandler retrieves tasks from a database and writes them as JSON to an HTTP
// response.
func tasks(w http.ResponseWriter, r *http.Request) {

	var tasks []*db.Task
	var err error
	search := r.FormValue("search")

	if len(search) == 0 {
		tasks, err = db.Tasks(limit)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else if date, err := time.Parse("02.01.2006", search); err == nil {
		dateStr := date.Format(layout)
		tasks, err = db.TasksDate(limit, dateStr)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		tasks, err = db.TasksSearch(limit, search)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if !idIsValid(id) {
		writeError(w, errNotValidID, http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJson(w, empty)

}
func taskDone(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if !idIsValid(id) {
		writeError(w, errNotValidID, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if len(task.Repeat) == 0 {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJson(w, empty)
		return
	}
	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err = db.UpdateDate(task.Date, id); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJson(w, empty)
}
func idIsValid(id string) bool {
	if len(id) == 0 {
		return false
	}
	if _, err := strconv.Atoi(id); err != nil {
		return false
	}
	return true
}
func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(layout)
		return nil
	}
	temp, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}
	if afterNow(now, temp) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(layout)
		} else {
			task.Date, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
