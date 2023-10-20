# Introduction

![Framework](promotions/banner.png)

![License](https://img.shields.io/github/license/ByPikod/go-crypto.svg?style=for-the-badge)
![Maintaned](https://img.shields.io/badge/Maintained%3F-yes-green.svg?style=for-the-badge)
![Commits](https://img.shields.io/github/commits-since/ByPikod/go-crypto/latest.svg?style=for-the-badge)
![Forks](https://img.shields.io/github/forks/ByPikod/go-crypto.svg?style=for-the-badge)
![Stars](https://img.shields.io/github/stars/ByPikod/go-crypto.svg?style=for-the-badge)
![Watchers](https://img.shields.io/github/watchers/ByPikod/go-crypto.svg?style=for-the-badge)

This project is the back-end prototype of a crypto wallet application I developed with the Go programming language during my internship at [ProxoLab](https://www.proxolab.com/) company.

I would also like to mention that I learned the Go programming language alongside this project. Additionally, this project allowed me to learn various DevOps concepts.

When I started writing the project, I had searched for many such examples, but the examples I found always had shortcomings. I am sure that there are many aspects of this project that can be improved. However, I believe it is an excellent example project because I designed it with this year's most commonly used technologies and best practices.

## Table of Contents

- [Introduction](#introduction)
  - [Table of Contents](#table-of-contents)
  - [To-do List](#to-do-list)
  - [Developer Tips](#developer-tips)
- [Getting Started](#getting-started)
  - [Compose Setup](#compose-setup)
    - [Requirements](#requirements)
    - [Installation](#installation)
  - [Kubernetes Setup](#kubernetes-setup)
    - [Kubectl](#kubectl)
    - [Monitoring Kubernetes](#monitoring-kubernetes)
      - [Lens](#lens)
      - [Docker Desktop](#docker-desktop)
- [Project Design](#project-design)
  - [Technologies](#technologies)
  - [Ports](#ports)
  - [Folder Structure](#folder-structure)
  - [Models](#models)
- [Crypto API](#crypto-api)
- [Auto Documentation](#auto-documentation)
  - [Swagger](#-swagger)
  - [Swaggo](#-swaggo)
- [Monitoring](#monitoring)
  - [Prometheus](#-prometheus)
  - [Grafana](#-grafana)
- [Logging](#logging)
  - [Traditional Logging](#traditional-logging)
  - [Loki](#-loki)
- [Load Test](#load-test)
  - [Using K6](#-using-k6)
  - [Monitoring Test Results](#monitoring-test-results)
- [API](#api)
- [Websocket](#websocket)
- [Copyright](#copyright)

## To-do List
While working on this project, you can find below the expectations that the company I interned at had from me:

- [x] Implement the following endpoints:
  - [x] User registration, login and me endpoints.
  - [x] An endpoint for listing crypto exchanges
  - [x] Endpoints for selling & buying crypto.
  - [x] Create a websocket endpoint to live update crypto exchanges.
- [x] Respond HTTP requests with Fiber
- [x] Provide a better database interface with GORM.
- [x] Dockerize the project to avoid version conflicts.
- [ ] Use the following database relationships:
  - [x] HasMany,
  - [ ] HasOne,
  - [x] BelongsTo,
  - [ ] ManyToMany
- [x] Use JWT to authenticate user.
- [x] Use GORM to migrate models.
- [x] Documentize the API with Swagger
- [ ] Test the app
  - [ ] Unit tests.
  - [x] Mock repository layer for unit tests.
- [x] Monitor the app
  - [x] Monitoring with Prometheus and Grafana
  - [x] Profile application with load test
  - [x] Monitor logs with Loki and Grafana.
- [ ] Create a micro service
  - [ ] Create a mail verification system for "Go Crypto".
  - [ ] Create a mail sender micro service named "Notifier" to send verification mails.
  - [ ] Connect "Notifier" and "Go Crypto" with **Kafka**
  - [ ] Use **Apache ZooKeeper** with **Kafka**

## Developer Tips

- File **"dockerfile.dev"** is in use for development purposes. You can change **"build"** property of **"gocrypto"** container from **docker-compose.yml**
- You can set **"HOST"** environment variable to listen to a specific domain with Fiber.
- Credentials for Grafana:
  - Username: `admin`
  - Password: `root`

# Getting Started

![Docker](promotions/docker-banner.png)

As you probably know, running a single Go file is not so hard. But when it comes to running a project with multiple dependencies, it gets a little bit complicated. And **Docker** comes to the rescue in this situation. **You can run the project with a single command using Docker**. You also need Docker for deployment purposes.

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers. Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

Docker can get your project up and running in a few minutes. But that's not enough for a project like this. You need to monitor your project, test it, and more. And for that you need to install some other tools like **Prometheus**, **Grafana**, **K6**, **Postgres**.

- [**Compose Setup:**](#compose-setup) Used for development purposes.
- [**Kubernetes Setup:**](#kubernetes-setup) Used for production purposes.

> [!NOTE]  
> **Both ways uses Docker Engine to run the project. You can install Docker Engine from [here](https://docs.docker.com/engine/install/).**

## Compose Setup

Compose Setup uses Docker Compose to run the project. Which is a tool designed to run multi-container applications. You need to configure a file named ["docker-compose.yml"](docker-compose.yml) to install all the containers needed and run the project. You can Click [here](https://docs.docker.com/compose/) to learn more about Docker Compose.

Go Crypto runs with the following containers:


| Container   | Description                                  |
|-------------|----------------------------------------------|
| Go-Crypto   | Runs the main project. Docker file can be found at project directory. |
| Notifier    | Runs the micro service that sends the mails. |
| K6          | Runs a script that simulates connection load to test servers. |
| Postgres    | Database to store user data.                 |
| Prometheus  | Database to store metric data.               |
| Loki        | Database to store logs.                      |
| Grafana     | Analytics visualizer interface.              |
| ZooKeeper   | Micro service manager.                       |
| Kafka       | Queue for micro services.                    |
| Kafka UI    | Interface for kafka. Used for development purposes |

These containers all work together to run the project.

### Requirements

You need to install the following tools to run the project with Docker Compose:

- [git-cli:](https://git-scm.com/downloads) Used to download the project to your device.
- [docker-compose:](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-20-04#step-1-installing-docker-compose) Used to quickly install required images (eg. PostgreSQL, Loki, Grafana).

### Installation

- Clone repository: `git clone https://www.github.com/ByPikod/go-crypto.git`
- Run docker-compose: `docker-compose up`
  - Optionally add "-d" arg to run at background: `docker-compose up -d`

## Kubernetes Setup

![Kubernetes](promotions/kubernetes.png)

Kubernetes is an open-source container orchestration platform used for automating the deployment, scaling, and management of containerized applications.

Kubernetes can automatically scale your applications based on demand, ensuring that you have the right amount of resources to handle your workload.

<img align="right" alt="Horizontal vs Verical scaling up" src="promotions/verticalvshorizontal.png" width="300px">

To understand what is Horizontally Scaling, first you should know what is Verical Scaling.

**Vertical scaling**, aka. scaling up, is a method used to increase the capacity and performance of a single server or machine by adding more resources to it (eg. RAM, CPU).

And **Horizontal scaling**,  aka. scaling out, is a method used to increase the capacity and performance of a system by adding more machines or nodes to a network or cluster.

You can scale out **Go-Crypto** by the following steps:

* Install Kubarnetes
* Apply configuration files: `kubectl apply -f .\kubernetes\`

And that's it, you are ready to go!

### Kubectl

Kubectl is a command-line tool that allows you to run commands against Kubernetes clusters. You can use kubectl to deploy applications, inspect and manage cluster resources, and view logs.

You can find some useful commands below:

- Listing
  - `kubectl get pods` - List all pods
  - `kubectl get services` - List all services
  - `kubectl get deployments` - List all deployments
  - `kubectl get nodes` - List all nodes
- Deleting Pods
  - `kubectl delete <type> <pod-name>` - Delete pods, services, deployments (Used to restart pods)
  - `kubectl delete <type> -l label=<label-value>` - Delete all pods, services, deployments, etc. by a label.
  - `kubectl delete <type> --all` - Delete all pods, services, deployments, etc of a type.
- Logs
  - `kubectl logs <pod-name>` - Get logs of a pod.
- Executing commands
  - `kubectl exec -it <pod-name> -- <command>` - Execute a command in a pod.

### Monitoring Kubernetes

#### Lens

[Lens](https://k8slens.dev/) is a powerful IDE for Kubernetes. It is a standalone application for MacOS, Windows and Linux operating systems. You can use Lens to monitor your Kubernetes cluster.

![Lens](promotions/lens.jpg)

#### Docker Desktop

Another tool you can use to monitor your Kubernetes cluster is [Docker Desktop](https://www.docker.com/products/docker-desktop). Docker Desktop is a tool for MacOS and Windows machines for the building and sharing of containerized applications and microservices.

![Docker Desktop](promotions/docker-desktop.png)

# Project Design

This project is a REST API project, and it has been built using the [Multitier Architecture](https://en.wikipedia.org/wiki/Multitier_architecture) design. The Multitier Architecture design is a REST API pattern consisting of layers named Controller (Presentation), Service (Business), and Repository (Persistence).

###s Technologies
You can find the technologies used in this project below:

* **Go:** Programming language.
* Libraries/Frameworks
    * **Gorn:** Database management
    * **Fiber:** HTTP library.
    * **JWT:** Auth standard.
* Documentation
    * **Swagger:** Rest API documentation.
    * **Swaggo:** Auto config generator for Swagger
* Database
    * **Postgres:** Database
    * **Prometheus:** Analytics collector.
* Tools
    * **Grafana:** Analytics visualizer.
    * **Air:** Live reload tool for go projects.
    * **Git:** Version control system
    * **Docker:** Containerization platform.

### Ports

The ports exposed by the project are:

| Name       | Port | Description                           |
|------------|------|---------------------------------------|
| Fiber      | 80   | Rest API itself.                      |
| Swagger    | 8080 | Auto generated documentation for API. |
| Grafana    | 3000 | Monitoring interface.                 |
| Postgres   | 5432 | Database                              |
| Prometheus | 9090 | Metric database.                      |
| Loki       | 3100 | Log database.                         |
| ZooKeeper  | 2181 | Micro services.                       |
| Kafka      | 9092 | Micro services.                       |

Some ports are exposed on development purposes.

### Folder Structure
This project follows common designs used in the back-end of web applications. Here is the project tree with comments explaining the modules, files, and their purposes:

```py
.
├── controllers # Endpoints (aka Presentation layer)
│   ├── exchanges.go
│   ├── user.go
│   └── wallet.go
├── core # Core components
│   ├── config.go # Retrieve environment variables
│   └── database.go # Initialize database connection
├── helpers # Utilities
│   ├── database.go # Database utilities
│   ├── errors.go # HTTP Errors (e.g 404, 403, 400)
│   ├── password.go # Password hashing, comparing
│   ├── token.go # JWT utilities
│   └── validate.go # Payload validations
├── log
│   ├── local.go # Logging functions for local logs.
│   └── log.go # Logging functions for Loki.
├── main.go
├── middleware
│   ├── auth.go # Authorization with JWT
│   ├── json.go # Adds header accepts "application/json"
│   ├── metrics.go # Monitoring
│   └── websocket.go # Return error if websocket request missing upgrade header.
├── models # Database & API models
│   ├── exchanges.go
│   ├── transaction.go
│   ├── user.go
│   └── wallet.go
├── repositories # Repository layer (aka Persistance)
│   ├── exchanges.go
│   ├── user.go
│   └── wallet.go
├── router # Routes
│   └── router.go
└── services # Service layer (aka Bussiness layer)
    ├── exchanges.go
    ├── user.go
    └── wallet.go
```

## Models

The module called "Models" is the boilerplate that represents data structures used in the background of web applications. Models are used to introduce raw data into the programming language being used. Here are the models for this project:

* **User:** Holds user data. Password encrypted with bcrypt.
    ```go
    type User struct {
        gorm.Model
        Name     string   `json:"name" gorm:"not null"`
        Lastname string   `json:"lastName" gorm:"not null"`
        Mail     string   `json:"mail" gorm:"index;not null;unique"`
        Password string   `json:"password" gorm:"not null"`
        Wallets  []Wallet `gorm:"foreignKey:UserID"` // HasMany
    }
    ```

* **Wallet:** User can have multiple wallets with different currencies for each one. These wallets can have transactions histories.
    ```go
    type Wallet struct {
        gorm.Model
        Currency    string        `json:"currency" gorm:"not null;index"`
        Balance     float64       `json:"balance" gorm:"default:0;not null"`
        UserID      uint          `json:"userID" gorm:"not null;index"`
        User        User          `gorm:"foreignKey:UserID"` // BelongsTo
        Transaction []Transaction `gorm:"foreignKey:WalletID"` // HasMany
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
        Wallet   Wallet  `gorm:"foreignKey:WalletID"` // BelongsTo
    }
    ```
* **Exchanges:** This model holds the exchange data received from API.
    ```go
    type ExchangeRates struct {
	    Currency string             `json:"currency"`
	    Rates    map[string]float64 `json:"rates"`
    }
    ```
    
# Crypto API

![Coinbase](/promotions/coinbase.jpg)

[Coinbase's](https://coinbase.com/) free to use public API (v2) is used to fetch crypto exchange data. Here you can find the endpoint below:

`https://api.coinbase.com/v2/exchange-rates?currency=BTC`

# Auto Documentation

### <img src="promotions/swagger.png" width="12"> Swagger
Swagger is a tool that helps developers to create document and test APIs (Application Programming Interfaces) for their software applications.

### <img src="promotions/swaggo.png" width="12"> Swaggo
Swagger needs a configuration file to create documents. This configuration is automatically generated by **Swaggo** from your comment lines.

You can read the documentation from "[localhost:8080](http://localhost:8080)" once you started docker.

# Monitoring

Prometheus and Grafana are tools used for monitoring and analyzing metrics of a web application. With Grafana and Prometheus, you can analyze a wide range of data, from sales metrics to resource utilization and more.

### <img src="promotions/prometheus-logo.png" width="14"> Prometheus

Prometheus is a software that collects "metric" data from the http servers by requesting a specific endpoint at target servers at intervals you specified. The metric data is stored chronologically by the Prometheus. Data can be accessed via web interface or Rest API that Prometheus provides.

![Prometheus web interface](promotions/prometheus.png)

### <img src="promotions/grafana.png" width="14"> Grafana

And Grafana is an open source analytics monitoring tool that provides bunch of visual components like (e.g charts, gauges). Grafana can have multiple data sources and Prometheus is one of them. Grafana can request to API of the Prometheus and visualize your chronologically stored metrics data.

I've configured Prometheus to gather **default Go Metrics** from **Go Fiber** and visualized some of those metrics in Grafana as can be seen in the picture below:

![monitoring](promotions/monitoring.png)

# Logging
This project uses two different methods for logging. One is for local logging with simple print functions and the other is for advanced logging for customers.

### Traditional Logging
I've used traditional local logging functions for the stuff that only concerns developers.

![traditional logging](promotions/traditional-logging.png)

### <img src="promotions/loki-logo.png" width="14"> Loki
Loki is an open-source log aggregation system developed by Grafana Labs. It is designed to help you collect, store, and query log data from various sources in a scalable and efficient manner. 

In this project, I used **Loki** for advanced logging. And I used [zap-loki](https://github.com/paul-milne/zap-loki) library for establishing the connection between Go and Loki.

Monitoring Loki data source in Grafana:

![logging](promotions/logging.png)

# Load Test

The term load testing is used in different ways in the professional software testing community. Load testing generally refers to the practice of modeling the expected usage of a software program by simulating multiple users accessing the program concurrently.
(https://en.wikipedia.org/wiki/Load_testing)

I've used "Grafana K6" to load test this project. K6 has a very simple interface.

### <img src="promotions/k6-logo.png" width="14"> Using K6
First, you should create a JavaScript file for K6, as mentioned, named ["scripts/loadtest.js"](scripts/loadtest.js).

This script, using the framework provided by K6, allows you to quickly send load requests to your web application. You can also perform checks on responses using the functions provided by K6.

K6 is originally designed to export metrics to a data source called InfluxDB. However, you can obtain output for Prometheus using the experimental [Prometheus Remote Write](https://k6.io/docs/results-output/real-time/prometheus-remote-write/) module.

### Monitoring Test Results
The output obtained from Prometheus can be visualized using the Grafana interface, as explained in the "Monitoring" section.

![loadtest monitoring](promotions/loadtest.png)

# API

* [Endpoints](#endpoints) /api/
* [User Endpoints](#user-endpoints) /api/user/
* [User Wallet Endpoints: ](#user-wallet-endpoints) /api/user/wallet/

### Endpoints

<!-- Currencies -->

<details>
<summary style="font-size: 1.5em;">
<code>GET</code> <code>/api/exchange-rates/</code>
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
If verification code is not provided, it will try to send a verification code to the mail address provided.
If the bool named "verificationSent" in response data is true, you must have received the mail address.

##### Parameters
| Name         | Type   | Description       |
|--------------|--------|-------------------|
| name         | string | First name        |
| lastName     | string | Last name         |
| mail         | string | Mail address      |
| password     | string | Password          |
| verification | string | Verification code |

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
<code>GET</code> <code>/api/user/me/</code>
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
<code>GET</code> <code>/api/user/wallet/balance/</code>
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