package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//Structs for user, merchant & transaction

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

type Merchant struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	Id        int    `json:"id"`
	From      int    `json:"from"`
	To        int    `json:"to"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type Page struct {
	Title string
}

// Secret token to encrypt/decrypt JWT tokens
var jwtTokenSecret = "abc123456def"

// Initialization
var Users []User
var Merchants []Merchant
var Transactions []Transaction
var tr_count = 2 // Keep track of transaction count, will be incremented for transaction IDs

// This function return all users in json format
func returnUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnUsers") // For logging purposes
	json.NewEncoder(w).Encode(Users)
}

// This function a specific user with id as a condition
func returnUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var uid int

	if i, err := strconv.Atoi(key); err == nil {
		uid = i
	}

	for _, user := range Users {
		if user.Id == uid {
			json.NewEncoder(w).Encode(user)
		}
	}
}

// This function return all merchants in json format
func returnMerchants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnMerchants")
	json.NewEncoder(w).Encode(Merchants)
}

// This function return all transactions in json format
func returnTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnTransactions")
	json.NewEncoder(w).Encode(Transactions)
}

// This display the login template as well as the login logic itself
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Login Page")

	conditionsMap := map[string]interface{}{}

	// Check if user already logged in
	username, _ := ExtractTokenUsername(r)

	if username != "" { // User already logged in
		conditionsMap["Username"] = username
		conditionsMap["LoginError"] = false

		http.Redirect(w, r, "/payment", http.StatusSeeOther)
	}

	// Verify username and password from POST value
	if r.FormValue("username") != "" && r.FormValue("password") != "" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		var pwd string

		// Check if username and password match from data
		for _, user := range Users {
			if user.Name == username {
				pwd = user.Password
			}
		}
		hashedPassword := []byte(pwd)

		// Use the bcrypt library to compare hash and password
		if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
			log.Println("Wrong credentials.")
			conditionsMap["LoginError"] = true
		} else {
			log.Println("Logged in :", username)
			conditionsMap["Username"] = username
			conditionsMap["LoginError"] = false

			// Create a new JSON Web Token
			tokenString, err := createToken(username)

			if err != nil {
				log.Println(err) // Quick handle
				os.Exit(1)
			}

			// Create the cookie for client - browser
			expirationTime := time.Now().Add(1 * time.Hour) // Set cookie expired after 1 hour

			cookie := &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/payment", http.StatusFound) // Redirect to payment page if login successful
		}

	}

	if err := templ.ExecuteTemplate(w, "login", &Page{Title: "Login"}); err != nil { // Stay at login page if login failed
		log.Println(err)
	}
}

func pay(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Pay Logic")

	// Check if user already logged in
	username, _ := ExtractTokenUsername(r)

	if username == "" { // User already logged in
		w.Write([]byte("Please login first."))
		return
	}

	var mid int
	var amt int
	var uid int

	// Payment logic
	if r.FormValue("merchant") != "" && r.FormValue("amount") != "" {

		// Convert values to int
		merchantID := r.FormValue("merchant")
		if i, err := strconv.Atoi(merchantID); err == nil {
			mid = i
		}
		amount := r.FormValue("amount")
		if i, err := strconv.Atoi(amount); err == nil {
			amt = i
		}

		// Decrease balance for designated user according to the amount input
		for i, user := range Users {
			if user.Name == username {
				if user.Balance < amt {
					fmt.Fprintf(w, "Insufficient balance!")
					return
				}

				uid = user.Id
				user.Balance = user.Balance - amt
				Users[i] = user
			}
		}

		// Increase balance for designated user according to the amount input
		for i, merchant := range Merchants {
			if merchant.Id == mid {

				merchant.Balance = merchant.Balance + amt
				Merchants[i] = merchant
			}
		}

		// Insert new transaction
		transaction := Transaction{Id: tr_count, From: uid, To: mid, Amount: amt, CreatedAt: time.Now().Format(time.RFC3339)}
		Transactions = append(Transactions, transaction)
		tr_count++ // Increase tr_count for insert (id_ purposes

		json.NewEncoder(w).Encode(transaction) // Returns the payment detail
	}
}

// Remove the "token" cookie for logout
func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Logout")
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("Cookie deleted. Logged out!\n"))
}

// Display the payment template
func payment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Payment Page")
	conditionsMap := map[string]interface{}{}

	// Check if user already logged in
	username, _ := ExtractTokenUsername(r)

	if username != "" {
		conditionsMap["Username"] = username
	} else {
		http.Redirect(w, r, "/login", http.StatusFound) // Redirect to login page if not yet logged in
	}

	// Execute payment template for logged in user
	if err := templ.ExecuteTemplate(w, "payment", conditionsMap); err != nil {
		log.Println(err)
	}
}

// Setup code for executing template
var templ = func() *template.Template {
	t := template.New("")
	err := filepath.Walk("./template/", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			fmt.Println(path)
			_, err = t.ParseFiles(path)
			if err != nil {
				fmt.Println(err)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}
	return t
}()

func createToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username                            // Embed username inside the token string
	claims["expired"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtTokenSecret))
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Generate Token API")

	vars := mux.Vars(r)
	username := vars["username"]
	tokenString, err := createToken(username)
	if err != nil {
		log.Println(err) // Quick handle
		os.Exit(1)
	}

	json.NewEncoder(w).Encode(tokenString)
}

func ExtractTokenUsername(r *http.Request) (string, error) {

	// Get our token string from Cookie
	biscuit, err := r.Cookie("token")

	var tokenString string
	if err != nil {
		tokenString = ""
	} else {
		tokenString = biscuit.Value
	}

	// Abort
	if tokenString == "" {
		return "", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtTokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username := fmt.Sprintf("%s", claims["username"]) // Convert to string
		if err != nil {
			return "", err
		}
		return username, nil
	}
	return "", nil
}

// Routes with Mux
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", payment)
	myRouter.HandleFunc("/login", login)
	myRouter.HandleFunc("/logout", logout)
	myRouter.HandleFunc("/payment", payment)
	myRouter.HandleFunc("/pay", pay).Methods("POST")
	myRouter.HandleFunc("/users", returnUsers)
	myRouter.HandleFunc("/merchants", returnMerchants)
	myRouter.HandleFunc("/transactions", returnTransactions)
	myRouter.HandleFunc("/user/{id}", returnUser)
	myRouter.HandleFunc("/token", generateToken).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Payments API in Go with Mux Routers")

	// Users initial data, Password: pass123
	Users = []User{
		User{Id: 1, Name: "kervin", Password: "$2b$10$knfJ.4rmOOdY5vDcOQZYeOPOsjO2Mwpl5KnPwN8j.4RIY5hekklSO", Balance: 100},
		User{Id: 2, Name: "admin", Password: "$2b$10$knfJ.4rmOOdY5vDcOQZYeOPOsjO2Mwpl5KnPwN8j.4RIY5hekklSO", Balance: 5000},
	}

	// Merchants initial data
	Merchants = []Merchant{
		Merchant{Id: 1, Name: "Merchant A", Balance: 1000},
		Merchant{Id: 2, Name: "Merchant B", Balance: 500},
		Merchant{Id: 3, Name: "Merchant C", Balance: 200},
	}

	// Transactions initial data
	Transactions = []Transaction{
		Transaction{Id: 1, From: 1, To: 1, Amount: 1000, CreatedAt: "2021-10-21 00:00:00"},
	}

	handleRequests()
}
