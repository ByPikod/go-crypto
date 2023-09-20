# Go Crypto
An intern project

### Goals
Prepare a prototype crypto wallet REST API and follow the rules below: 

* Respond HTTP requests with Fiber
* Provide a better database interface with GORM.
* Dockerize the project to avoid version conflicts.
* Use the following postgres normalization rules: HasMany, HasOne, BelongsTo, ManyToMany
* Use JWT for the authentication.
* Implement the GORN auto migrate.

### Technologies Used
* Go (Must)
* Gorn (Must)
* Fiber (Must)
* JWT (Must)
* Docker (Must)
* Postgres (Must)
* Air
* Postman

### Structure
```py
.
├── core/ 
│   ├── config.go   # Configuration
│   ├── database.go # Database initialize and access
│   └── server.go   # Start to listen HTTP requests
├── logger/
│   └── log.go      # Logging tools
├── middleware/
│   └── json.go     # Middleware for accepting JSON requests
├── models/
│   ├── base.go     # Database table base.
│   └── user.go     # User table
├── routes/
│   ├── http-errors/         # Routes for HTTP errors (404, 403 ...) 
│   ├── user/                # User endpoints
│   └── exchange-rates.go    # List exchange rates.
├── workers/    # Async tasks
└── main.go     # Main file
```

### Endpoints
* **/api/user/login:** Returns auth token if matching credentials provided.
```json
{
    "message": "OK!",
    "status": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.HBNfNTMv3Jd9Wf-m3v6buHgGLQL0Srl8zwGro8JHcO4"
}
```

* **/api/user/register:** Creates a new account if provided details are appropriate.
```json
{
    "message": "OK!",
    "status": true
}
```

* **/api/user/me:** Returns user information.

* **/api/user/logout:** Returns user session.

* **/api/user/balance:** Returns user balance.
```json
{
    "currency": "USD",
    "rates": {
        "00": 13.486176668914363,
        "1INCH": 3.9920159680638725,
        "AAVE": 0.0157022846824213,
        "ABT": 13.297872340425531,
        "ACH": 66.59119664380368,
        "ACS": 533.4613640607079,
        ...
    }
}
```

* **/api/user/wallet/deposit:** A money deposit endpoint. Virtual POS not implemented. It's just a prototype.
```json
{
    "newBalance": 64.05993807839195,
    "mesasge": "OK",
    "status": true
}
```

* **/api/user/wallet/withdraw:** A money withdraw endpoint.
```json
{
    "newBalance": 64.05993807839195,
    "mesasge": "OK",
    "status": true
}
```

* **/api/user/wallet/buy:** Performs a crypto purchase and returns success or failure depending on the result.
```json
{
    "Balance": { // New balance
        "BTC": 1.0053999999999998,
        "USD": 64.05993807839195
    },
    "sold_amount": 135.72802500006378,
    "sold_currency": "USD",
    "bought_amount": 0.005,
    "bought_currency": "BTC",
    "message": "OK!",
    "status": true // Success state
}
```

* **/api/user/wallet/sell:** Performs a crypto selling and returns success or failure depending on the result.


* **/api/currencies:** Lists the crypto currencies.

# Copyright
Copyright (c) 2023, [Yahya Batulu](https://www.yahyabatulu.com). Released under [MIT License](LICENSE)