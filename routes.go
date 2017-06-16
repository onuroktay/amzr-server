package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/onuroktay/amazon-reader/amzr-server/elasticsearch"
	"github.com/onuroktay/amazon-reader/amzr-server/item-data"
	"github.com/onuroktay/amazon-reader/amzr-server/user-data"
	"github.com/onuroktay/amazon-reader/amzr-server/util"
	"github.com/onuroktay/amazon-reader/analyse-fichier-json/step10"
	"os"
)

// SubHits contains a sub-data structure returns by elasticsearch

// Response contains answer send to client
type Response struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}
//toto
func routes() {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	r.HandleFunc("/", startPage).Methods("GET")
	r.PathPrefix("/app/").HandlerFunc(startPage).Methods("GET")
	//r.HandleFunc("/login", startPage).Methods("GET")
	//r.HandleFunc("/register", startPage).Methods("GET")

	r.HandleFunc("/login", checkLoginValidity).Methods("POST") // check login Validity
	r.HandleFunc("/logout", logout).Methods("GET")             // logout

	/* CRUD USERS */

	// Create new user in DB
	r.HandleFunc("/user", addUser).Methods("POST")

	// Get User information in DB
	//r.Handle("/user/{id}", checkAuth(ADMIN,
	//	http.HandlerFunc(getUser))).Methods("GET")

	// Get users information
	r.Handle("/users", checkAuth(OnurTPIUser.ADMIN,
		http.HandlerFunc(getUsers))).Methods("GET")

	// Preflight Request for users PUT and DELETE
	r.HandleFunc("/user/{id}/{role}", preflightRequest).Methods("OPTIONS") // preflight request

	// Update user
	r.Handle("/user/{id}/{role}", checkAuth(OnurTPIUser.ADMIN,
		http.HandlerFunc(updateUser))).Methods("PUT") // update rolevalue in db

	// Delete user
	r.HandleFunc("/user/{id}", preflightRequest).Methods("OPTIONS") // preflight request
	r.Handle("/user/{id}", checkAuth(OnurTPIUser.ADMIN,
		http.HandlerFunc(deleteUser))).Methods("DELETE") // update rolevalue in db

	/* RUD ITEMS */

	// Get Item information in DB
	r.Handle("/item/{id}", checkAuth(OnurTPIUser.USER,
		http.HandlerFunc(getItem))).Methods("GET")

	// Get Items information in DB
	r.Handle("/items", checkAuth(OnurTPIUser.USER,
		http.HandlerFunc(getItems))).Methods("POST") // get items in db

	// Preflight Request for items PUT and DELETE
	r.HandleFunc("/item/{id}", preflightRequest).Methods("OPTIONS") // preflight request

	// Update item
	r.Handle("/item/{id}", checkAuth(OnurTPIUser.EDITOR,
		http.HandlerFunc(updateItem))).Methods("PUT") // update item

	// Delete item
	r.Handle("/item/{id}", checkAuth(OnurTPIUser.EDITOR,
		http.HandlerFunc(deleteItem))).Methods("DELETE") // delete item
	// <----

	r.Handle("/import", checkAuth(OnurTPIUser.ADMIN,
		http.HandlerFunc(executeImport))).Methods("POST") // import data from json file

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(path))) // to get resources (like http, js, css, png, ...)
	http.Handle("/", r)
}

// This is the entry point of the client
func startPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	http.ServeFile(w, r, path + "index.html")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.accesser.GetUsers() // We get all user in DB
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, users, true)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idUser := param["id"]

	account, err := database.accesser.GetAccountByIDInDB(idUser) // We check if the user is in DB
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	//Successful
	writeResponse(w, r, account.User, true)
}

func addUser(w http.ResponseWriter, r *http.Request) {

	var credential *OnurTPIUser.CredentialsClient

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 2000000))

	// if the error is different from null
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(err.Error())
		w.Write(resp)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Print(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(err.Error())
		w.Write(resp)
		return
	}

	// we do unmarshal login of the client
	err = json.Unmarshal(body, &credential)
	if err != nil {
		log.Print(err.Error())
		// en cas d'échec
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(err.Error())
		w.Write(resp)
		return
	}

	found, _ := database.accesser.CheckIfUserExistsInDB(credential) // We check if the user is in DB
	if found == true {
		writeResponse(w, r, []byte("Sorry, username already exists!"), false)
		return
	}

	account := credential.CreateAccount()

	// Save the login in the database
	idUser, err := database.accesser.CreateAccount(account)
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, idUser, true)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idUser := param["id"]
	userRole, err := strconv.Atoi(param["role"])

	if err != nil {
		writeResponse(w, r, err, false)
		return
	}



	// save the login in the database
	err = database.accesser.UpdateRole(idUser, userRole)

	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, idUser, true)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idUser := param["id"]

	// save the login in the database
	err := database.accesser.DeleteAccount(idUser)

	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, idUser, true)
}

func checkLoginValidity(w http.ResponseWriter, r *http.Request) {
	var credential *OnurTPIUser.CredentialsClient // data from Client
	var maxSize int64 = 1024
	var tokenString string

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxSize)) //Read Body with size 1024
	if err != nil {
		sendMaxSizeExceeded(w, maxSize)
		return // If we arrive until then on does not execute the rest of code
	}
	if err := r.Body.Close(); err != nil {
		writeResponse(w, r, err, false)
		return // If we arrive until then on does not execute the rest of code
	}

	// We read the body
	if err := json.Unmarshal(body, &credential); err != nil {
		writeResponse(w, r, err, false)
		return
	}

	account, err := database.accesser.GetAccountByUserNameInDB(credential) // We check if the user is in DB
	if account == nil || err != nil {
		writeResponse(w, r, "Wrong login", false)
		return
	}

	// Comparing the password
	if !OnurTPIUser.CompareHashAndPassword(account.Password, credential.Password) {
		writeResponse(w, r, "Wrong login", false)
		return
	}

	// Create Token, with session data
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"IdUser": account.ID,
		"Expire": time.Now().Add(24 * time.Hour * 32), // 32 days
	})

	// Sign and get the complete encoded token as a string using the secret key
	tokenString, err = token.SignedString(util.GetEncryptionKey())
	if err != nil {
		fmt.Println(tokenString, err)
	}

	userRole := 0
	expire := time.Now()
	switch userRole {
	case 1:
		expire = time.Now().AddDate(0, 1, 0)
	case 2:
		expire = time.Now().AddDate(0, 0, 7)
	default:
		expire = time.Now().AddDate(0, 0, 1)

	}

	//expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:     "onurAuth",
		Value:    tokenString,
		Path:     "/",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	writeResponse(w, r, &OnurTPIUser.User{
		UserName:  account.UserName,
		RoleValue: account.RoleValue,
	}, true)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idItem := param["id"]

	item, err := database.accesser.GetItemByIDInDB(idItem) // We check if the item is in DB
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	//Successful
	writeResponse(w, r, item, true)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idItem := param["id"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 2000000))
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	var item *OnurTPIItem.Item

	err = json.Unmarshal(body, &item)
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// save the item in the database
	err = database.accesser.UpdateItem(idItem, item)
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, item, true)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idItem := param["id"]

	// save the login in the database
	err := database.accesser.DeleteItem(idItem)

	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, idItem, true)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var searchFromClient *OnurTPIES.Search

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 2000000))

	// si l'erreur différent du null
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	err = json.Unmarshal(body, &searchFromClient)
	if err != nil {
		fmt.Println(string(body))
		writeResponse(w, r, err, false)
		return
	}

	resp, err := database.accesser.GetItems(searchFromClient) // We get all item in DB
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	// Successful
	writeResponse(w, r, resp, true)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "onurAuth",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		//Secure:   true,
		//HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	//Successful
	writeResponse(w, r, cookie, true)
}

func executeImport(w http.ResponseWriter, r *http.Request) {

	// Read json content
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, r, err, false)
		return
	}

	fileName := "../AmzR-import-cli/" + util.GetUUID() + ".json"

	// Open file for writing
	fileW, err := os.Create(fileName)
	// Save content in file

	if err != nil {
		log.Fatal(err)
	}
	defer fileW.Close()

	n3, err := fileW.WriteString(string(body))
	fmt.Printf("wrote %d bytes\n", n3)
	fileW.Sync()

	// Import
	err = OnurTPIjsonReader.ImportJSON(fileName)
	if err != nil {
		fmt.Println(err)
	}

	// Delete file
	err = os.Remove(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Successful
	writeResponse(w, r, true, true)
}

// setHeader adds header to the http response
func setHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Methods", "DELETE,PUT,GET,POST")
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
}

// sendMaxSizeExceeded defines a maximal size for the body
func sendMaxSizeExceeded(w http.ResponseWriter, maxSize int64) {
	errMessage := "Sorry, the max size (" + CastInt64ToString(maxSize) + ") is exceeded"

	w.WriteHeader(http.StatusNotAcceptable)
	w.Write([]byte(errMessage))
}

// CastInt64ToString convert int64 to string
func CastInt64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func writeResponse(w http.ResponseWriter, r *http.Request, data interface{}, success bool) {
	setHeader(w, r)

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(&Response{
		Data:    data,
		Success: success,
	})

	w.Write(resp)
}

func preflightRequest(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)

	jsonResponse, _ := json.Marshal(true)
	w.Write(jsonResponse)
}
