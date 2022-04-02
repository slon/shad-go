//go:build !change

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"gitlab.com/slon/shad-go/coverme/models"
	"gitlab.com/slon/shad-go/coverme/utils"
)

type App struct {
	router *mux.Router
	db     models.Storage
}

func New(db models.Storage) *App {
	return &App{db: db}
}

func (app *App) Start(port int) {
	app.initRoutes()
	app.run(fmt.Sprintf(":%d", port))
}

func (app *App) initRoutes() {
	app.router = mux.NewRouter()
	app.router.HandleFunc("/", app.status).Methods("GET")
	app.router.HandleFunc("/todo", app.list).Methods("GET")
	app.router.HandleFunc("/todo/{id:[0-9]+}", app.getTodo).Methods("GET")
	app.router.HandleFunc("/todo/{id:[0-9]+}/finish", app.finishTodo).Methods("POST")
	app.router.HandleFunc("/todo/create", app.addTodo).Methods("POST")
}

func (app *App) run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stderr, app.router)
	_ = http.ListenAndServe(addr, loggedRouter)
}

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	todos, err := app.db.GetAll()
	if err != nil {
		utils.ServerError(w)
		return
	}

	_ = utils.RespondJSON(w, http.StatusOK, todos)
}

func (app *App) addTodo(w http.ResponseWriter, r *http.Request) {
	req := &models.AddRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		utils.BadRequest(w, "payload is required")
		return
	}
	defer func() { _ = r.Body.Close() }()

	if req.Title == "" {
		utils.BadRequest(w, "title is required")
		return
	}

	todo, err := app.db.AddTodo(req.Title, req.Content)
	if err != nil {
		utils.ServerError(w)
		return
	}

	_ = utils.RespondJSON(w, http.StatusCreated, todo)
}

func (app *App) getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
		return
	}

	todo, err := app.db.GetTodo(models.ID(id))
	if err != nil {
		utils.ServerError(w)
		return
	}

	_ = utils.RespondJSON(w, http.StatusOK, todo)
}

func (app *App) finishTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.BadRequest(w, "ID must be an int")
		return
	}

	if err := app.db.FinishTodo(models.ID(id)); err != nil {
		utils.ServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	_ = utils.RespondJSON(w, http.StatusOK, "API is up and working!")
}
