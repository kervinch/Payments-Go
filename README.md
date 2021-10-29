# Payments-Go

A simple payment simulation program in Go.

The functionality consist of login, payment and logout. <br>
The models consist of users, merchants and transactions.

#### User Model:

        type User struct {
	        Id       int    `json:"id"`
	        Name     string `json:"name"`
	        Password string `json:"password"`
	        Balance  int    `json:"balance"`
        }
        
#### Merchant Model:

        type User struct {
	        Id       int    `json:"id"`
	        Name     string `json:"name"`
	        Balance  int    `json:"balance"`
        }
        
#### Transaction Model:

        type Transaction struct {
          Id        int    `json:"id"`
          From      int    `json:"from"`
          To        int    `json:"to"`
          Amount    int    `json:"amount"`
          CreatedAt string `json:"created_at"`
        }
        
<hr>
The APIs are as follows: <br>

## Auth
### Login
[POST] /login <br>
#### Parameters: <br>
username string <br>
password string <br>
#### Response: <br>
-- set login cookie & redirects to payment.html --

<br>

### Logout
[GET] /logout <br>
#### Parameters: none <br>
#### Response: <br>
-- remove login cookie --

<br>

## Payment
### Pay
[POST] /pay <br>
#### Parameters: <br>
merchant int <br>
amount int <br>
#### Response: <br>
        {
            "id": 2,
            "from": 1,
            "to": 3,
            "amount": 40,
            "created_at": "2021-10-29T10:22:50+07:00"
        }
        
<br>

### Payment Page
[GET] /payment <br>
#### Parameters: none <br>
#### Response: <br>
-- execute payment template, payment.html --

<br>

## Users
### Get users
[GET] /users <br>
#### Parameters: none <br>
#### Response: <br>
        [
            {
                "id": 1,
                "name": "kervin",
                "password": "$2b$10$knfJ.4rmOOdY5vDcOQZYeOPOsjO2Mwpl5KnPwN8j.4RIY5hekklSO",
                "balance": 90
            },
            {
                "id": 2,
                "name": "admin",
                "password": "$2b$10$knfJ.4rmOOdY5vDcOQZYeOPOsjO2Mwpl5KnPwN8j.4RIY5hekklSO",
                "balance": 5000
            }
        ]
        
<br>

### Get user
[GET] /user/{id} <br>
#### Parameters: <br>
id int <br>
#### Response: <br>
        {
            "id": 2,
            "name": "admin",
            "password": "$2b$10$knfJ.4rmOOdY5vDcOQZYeOPOsjO2Mwpl5KnPwN8j.4RIY5hekklSO",
            "balance": 5000
        }
        
<br>

## Merchants
### Get merchants
[GET] /merchants <br>
#### Parameters: none <br>
#### Response: <br>
        [
            {
                "id": 1,
                "name": "Merchant A",
                "balance": 1000
            },
            {
                "id": 2,
                "name": "Merchant B",
                "balance": 500
            },
            {
                "id": 3,
                "name": "Merchant C",
                "balance": 240
            }
        ]
        
<br>

## Transactions
### Get transactions
[GET] /transactions <br>
#### Parameters: none <br>
#### Response: <br>
        [
            {
                "id": 1,
                "from": 1,
                "to": 1,
                "amount": 1000,
                "created_at": "2021-10-21 00:00:00"
            },
            {
                "id": 2,
                "from": 0,
                "to": 3,
                "amount": 40,
                "created_at": "2021-10-29T10:22:50+07:00"
            }
        ]
        
<br>

## How to use

Run the command: <b>go run main.go </b> to run the program. <br>
You will be prompted to login before you can access the main page.  <br><br>
Use this credential: <br>
username: admin <br>
password: pass123 <br>

You will be redirected to the payments page. Here you can select which merchant you want to transfer money to, choices are Merchant A, B and C. <br>
After you click pay, balance from user admin will be deducted and balance of the designated merchant will increase. <br>
You can open http://127.0.0.1:10000/transactions to see the list of transactions. <br>

Core routes: <br>
http://127.0.0.1:10000/login <br>
http://127.0.0.1:10000/payment <br>
http://127.0.0.1:10000/transactions <br>
http://127.0.0.1:10000/users <br>
http://127.0.0.1:10000/merchants <br>
http://127.0.0.1:10000/logout <br>
