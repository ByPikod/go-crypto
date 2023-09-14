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

* **/api/user/register:** Creates a new account if provided details are appropriate.

* **/api/user/me:** Returns user information.

* **/api/user/logout:** Returns user session.

* **/api/user/balance:** Returns user balance.

* **/api/user/buy:** Performs a crypto purchase and returns success or failure depending on the result.

* **/api/user/sell:** Performs a crypto selling and returns success or failure depending on the result.

* **/api/currencies:** Lists the crypto currencies.

# Copyright
Copyright (c) 2023, [Yahya Batulu](https://www.yahyabatulu.com). Released under [MIT License](LICENSE)