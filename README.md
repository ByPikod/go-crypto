# Introduction
This is a prototype back-end of a crypto application that I developed for my internship. 

#### Table of Contents
- [Introduction](#introduction)
    - [Goal of the Project](#goal-of-the-project)
    - [Technologies Used](#technologies-used)
- [Project Design](#project-design)
    - [Folder Structure](#folder-structure)
    - [Models](#models)
- [API Documentation](#api)
- [Copyright](#copyright)

#### To-do List
Prepare a prototype crypto wallet REST API and follow the rules below: 

* [x] Respond HTTP requests with Fiber
* [x] Provide a better database interface with GORM.
* [x] Dockerize the project to avoid version conflicts.
* [ ] Use the following postgres normalization rules: 
    * [x] HasMany, 
    * [ ] HasOne, 
    * [ ] BelongsTo,
    * [ ] ManyToMany
* [x] Use JWT for the authentication.
* [x] Implement the GORN auto migrate.
* [x] Documentize API with Swagger
* [x] Unit tests (at least one unit)

# Installation

**Requirements:**
* [docker-compose](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-20-04#step-1-installing-docker-compose)

**Installation:**

* Clone repository
    > git clone https://www.github.com/ByPikod/go-crypto.git
* Run docker-compose:
    > docker-compose up

    Optionally add "-d" arg to run at background

    > docker-compose up -d    

# Project Design

#### Dependencies
![dependencies](promotions/dependencies.png)

#### Technologies Used
* Must:
    * Go
    * Gorn
    * Fiber
    * JWT
    * Docker
    * Postgres
    * Swagger
    * Swaggo
* Optional:
    * Air
    * Postman

## Folder Structure
```py
.
├── core # Core components
│   ├── config.go # Environment variables
│   └── database.go # Global DB utilities and DB connection
├── helpers # Utilities
│   ├── logger.go # Logging utilities
│   └── password.go # Password hashing, validating etc.
├── main.go # Launcher
├── middleware # Middle functions between framework and routes.
│   ├── auth.go # Checks if user is logged in
│   └── json.go # Adds "accept json" header
│   └── websocket.go # Returns "upgraded needed" if request is not applicable with ws.
├── models # Database tables
│   ├── transaction.go
│   ├── transaction_service.go
│   ├── user.go
│   ├── user_service.go
│   ├── wallet.go
│   └── wallet_service.go
├── routes
│   ├── errors.go # Error codes such as 400, 401
│   ├── exchanges.go # Exchange list end point
│   ├── user.go # User related endpoints
│   └── wallet.go # User wallet related endpoints
└── workers # Each thread is a different worker.
    ├── exchanges.go # Fetchs exchange rates from a public API with an interval.
    └── router # HTTP Server
        └── router.go
```

## Models

* **User:** Holds user data. Password encrypted with bcrypt.
    ```go
    type User struct {
        gorm.Model
        Name     string   `json:"name" gorm:"not null"`
        Lastname string   `json:"lastName" gorm:"not null"`
        Mail     string   `json:"mail" gorm:"index;not null;unique"`
        Password string   `json:"password" gorm:"not null"`
        Wallets  []Wallet `gorm:"foreignKey:UserID"` // Has Many
    }
    ```

* **Wallet:** User can have multiple wallets with different currencies for each one. These wallets can have transactions histories.
    ```go
    type Wallet struct {
        gorm.Model
        Currency    string        `json:"currency" gorm:"not null;index"`
        Balance     float64       `json:"balance" gorm:"default:0;not null"`
        UserID      uint          `json:"userID" gorm:"not null;index"`
        Transaction []Transaction `gorm:"foreignKey:WalletID"` // Has Many
    }
    ```

* **Transaction:** Transaction history holds the history of transactions as the name describes.
    ```go
    type Transaction struct {
        gorm.Model
        Type     int8    `json:"type" gorm:"not null"`
        Change   float64 `json:"change" gorm:"not null"`
        Balance  float64 `json:"balance" gorm:"not null"`
        WalletID uint    `json:"walletID" gorm:"not null;index"`
    }
    ```

# API

* [Endpoints](#endpoints) /api/
* [User Endpoints](#user-endpoints) /api/user/
* [User Wallet Endpoints: ](#user-wallet-endpoints) /api/user/wallet/

### Endpoints

<!-- Currencies -->

<details>
<summary style="font-size: 1.5em;">
<code>GET</code> <code>/api/currencies/</code>
</summary>

##### Description    
Lists the crypto currency exchange rates.

##### Response
    
```json
{
    "currency": "USD",
    "rates": {
        "00": 13.651877133105803,
        "1INCH": 3.898635477582846,
        "AAVE": 0.0159936025589764,
        "ABT": 13.708019191226867,
        "ACH": 64.1148938898506,
        ...
}
```
</details>

### User Endpoints

<!-- Login -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/login/</code>
</summary>

##### Description    
Returns auth token if matching credentials provided.

##### Parameters
| Name     | Type   | Description  |
|----------|--------|--------------|
| mail     | string | Mail address |
| password | string | Password     |

##### Response
    
```json
{
    "status": true,
    "message": "OK!",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.HBNfNTMv3Jd9Wf-m3v6buHgGLQL0Srl8zwGro8JHcO4"
}
```
</details>

<!-- Register -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/register/</code>
</summary>

##### Description    
Creates a new account if provided details are appropriate.

##### Parameters
| Name     | Type   | Description  |
|----------|--------|--------------|
| name     | string | First name   |
| lastName | string | Last name    |
| mail     | string | Mail address |
| password | string | Password     |

##### Response
    
```json
{
    "status": true,
    "message": "OK!",
}
```
</details>

<!-- Me -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/me/</code>
</summary>

##### Description    
Returns user information.

##### Parameters
Auth required

##### Response
    
```json
{
    "id": 1,
    "lastname": "Batulu",
    "mail": "admin@yahyabatulu.com",
    "name": "Yahya"
}
```
</details>

### User Wallet Endpoints

<!-- Balance -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/wallet/balance/</code>
</summary>

##### Description    
Returns user balance.

##### Parameters
Auth required

##### Response
    
```json
{
    "BTC": 40990.47869000058,
    "USD": 995270.5766880848
}
```
</details>

<!-- Deposit -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/wallet/deposit/</code>
</summary>

##### Description    
A money deposit endpoint. Virtual POS not implemented. It's just a prototype.

##### Parameters
Auth required

| Name   | Type  | Description                 |
|--------|-------|-----------------------------|
| amount | float | Amount of money to deposit. |

##### Response
    
```json
{
    "status": true,
    "newBalance": 64.05993807839195,
    "message": "OK",
}
```
</details>

<!-- Withdraw -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/wallet/withdraw/</code>
</summary>

##### Description    
A money withdraw endpoint.

##### Parameters
Auth required

| Name   | Type  | Description                  |
|--------|-------|------------------------------|
| amount | float | Amount of money to withdraw. |

##### Response
    
```json
{
    "status": true,
    "newBalance": 64.05993807839195,
    "message": "OK",
}
```
</details>

<!-- Buy -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/wallet/buy/</code>
</summary>

##### Description    
Performs a crypto purchase and returns success or failure depending on the result.

##### Parameters
Auth required

| Name     | Type   | Description              |
|----------|--------|--------------------------|
| amount   | float  | Amount of crypto to buy. |
| currency | string | Currency to buy          |

##### Response
    
```json
{
    "message": "OK!",
    "status": true, // Success state
    "Balance": { // New balance
        "BTC": 1.0053999999999998,
        "USD": 64.05993807839195
    },
    "sold_amount": 135.72802500006378,
    "sold_currency": "USD",
    "bought_amount": 0.005,
    "bought_currency": "BTC",
}
```
</details>

<!-- Buy -->

<details>
<summary style="font-size: 1.5em;">
<code>POST</code> <code>/api/user/wallet/sell/</code>
</summary>

##### Description    
Performs a crypto selling and returns success or failure depending on the result.

##### Parameters
Auth required

| Name     | Type   | Description               |
|----------|--------|---------------------------|
| amount   | float  | Amount of crypto to sell. |
| currency | string | Currency to sell.         |

##### Response
    
```json
{
    "message": "OK!",
    "status": true,
    "Balance": {
        "BTC": 995270.5766880848,
        "USD": 40990.47869000058
    },
    "bought_amount": 135.17072499994828,
    "bought_currency": "USD",
    "sold_amount": 0.005,
    "sold_currency": "BTC",
}
```
</details>

# Websocket

<!-- Balance -->

<details>
<summary style="font-size: 1.5em;">
<code>WS</code> <code>/ws/exchange-rates</code>
</summary>

##### Description    
Returns exchange-rates as it changed.

##### Response
    
```json
{
    "currency": "USD",
    "rates": {
        "00": 13.651877133105803,
        "1INCH": 3.898635477582846,
        "AAVE": 0.0159936025589764,
        "ABT": 13.708019191226867,
        "ACH": 64.1148938898506,
        ...
}
```
</details>

# Copyright
This project is licensed under the terms of the MIT License.

You are free to use this project in compliance with the MIT License. If you decide to use, modify, or redistribute this software, you must include a copy of the original license and copyright notice in all copies or substantial portions of the software.

For more information about the MIT License, visit: [MIT License](LICENSE).