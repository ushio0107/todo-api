# Todo List API
This is a Todo List API with CRUD functionality implemented using Go and MongoDB. The API provides the following features:

- Login: Users can authenticate and obtain a JWT token to access the API.
- Todo List Operations: Users can create, read, update, and delete Todo items.

## Prerequisites
Before getting started with the API, make sure you have the following software and services:

- Go Programming Language: https://golang.org/doc/install
- MongoDB Database: https://www.mongodb.com/

## Installation and Setup
1. Download or clone the source code of this project.

```git clone https://github.com/ushio0107/todo-api.git```
2. Create a new database in MongoDB.
3. Edit the config.go file and update the database connection URL with your MongoDB configuration.

```
const mongoDBUrl = "mongodb://localhost:27017"
const mongoDBName = "your_database_name"
```
4.Build and run the API.
```
go build
./todo-api
```
5. The API will be running locally at http://localhost:3000/v1.

## Using the API
Before using the API, obtain a JWT token. Send a POST request to the `/login` endpoint with the correct username and password.

```
POST /login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```
Upon successful login, you will receive a JWT token.

### Create a Todo
Send a POST request to the `/v1/todos` endpoint with the details of the Todo item.

```
POST /v1/todos
Authorization: Bearer <token>
Content-Type: application/json

{
  "task": "Complete the project",
  "completed": false
}
```

### Read All Todos
Send a GET request to the `/v1/todos` endpoint to retrieve all Todo items.

```
GET /v1/todos
Authorization: Bearer <token>
```

### Read a Single Todo
Send a GET request to the `/v1/todos/{id}` endpoint to retrieve a specific Todo item by its ID.

```
GET /v1/todos/{id}
Authorization: Bearer <token>
```

### Update a Todo
Send a PUT request to the `/v1/todos/{id}` endpoint to update a specific Todo item by its ID.

```
PUT /v1/todos/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "task": "Update the project",
  "completed": true
}
```

### Delete a Todo
Send a DELETE request to the `/v1/todos/{id}` endpoint to delete a specific Todo item by its ID.

```
DELETE /v1/todos/{id}
Authorization: Bearer <token>
```
