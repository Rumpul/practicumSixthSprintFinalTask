package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, `"Bad Request"`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка при записи ответа:", err)
		http.Error(w, `"Bad Request"`, http.StatusInternalServerError)
		return
	}

}

func getOneTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		http.Error(w, `"Задачи с таким ID не существует"`, http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка при записи ответа:", err)
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(bodyBytes, &task); err != nil {
		http.Error(w, `"Bad Request"`, http.StatusBadRequest)
		return
	}
	_, ok := tasks[task.ID]
	if ok {
		http.Error(w, "Задача с таким ID уже cуществует", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.WriteHeader(http.StatusCreated)
}

func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задачи с таким ID не существует", http.StatusBadRequest)
		return
	}

	delete(tasks, task.ID)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Get("/tasks/{id}", getOneTask)
	r.Post("/tasks", addTask)
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
