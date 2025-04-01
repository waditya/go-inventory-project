package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Sturct to store Router and DB information
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Method Initialise of struct to create connection string for database connection
func (app *App) Initialise() error {
	err := godotenv.Load()
	checkError(err)

	// Fetch the environment variables loaded by GoDotEnv in the above step
	DbUser := os.Getenv("DATABASE_USER")
	DbPassword := os.Getenv("DATABASE_PASSWORD")
	DbName := os.Getenv("DATABASE_NAME")

	fmt.Println(DbUser)
	fmt.Println(DbPassword)
	fmt.Println(DbName)

	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	// connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", "root", "adminpassword", "inventory")
	// connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", "root", "adminpassword", "learning")

	app.DB, err = sql.Open("mysql", connectionString)
	fmt.Println("Attemting to Open a connection with database")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// defer app.DB.Close()

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()

	return nil
}

func (app *App) Run(address string) {
	fmt.Println("Server starting to listen on port 10000")
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	// Convert payload to json
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func (app *App) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Store the map  of variables in a request determined by the Mux router in vars.
	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	p := product{ID: key}

	err = p.getProduct(app.DB) // Pass the pointer referencing the DB connection pool to the getProduct method

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "product not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponse(w, http.StatusOK, p)

}

func (app *App) createProductHandler(w http.ResponseWriter, r *http.Request) {

	var p product

	// Decode the data from request JSON payload

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request id")
		return
	}

	err = p.createProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) updateProductHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Store the map  of variables in a request determined by the Mux router in vars.
	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var p product

	// Decode the data from request JSON payload

	err = json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	p.ID = key

	err = p.updateProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProductsHandler).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProductHandler).Methods("GET")
	app.Router.HandleFunc("/product/", app.createProductHandler).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProductHandler).Methods("PUT")
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
