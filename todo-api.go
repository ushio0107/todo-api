package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Task      string             `json:"task,omitempty" bson:"task,omitempty"`
	Completed bool               `json:"completed,omitempty" bson:"completed,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

var collection *mongo.Collection
var signingKey = []byte("secret")

func ValidateToken(next http.Handler) http.Handler {
	log.Println("Validating token...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Logining...")
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Username != "admin" || user.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var todos []Todo
	for cursor.Next(ctx) {
		var todo Todo
		err := cursor.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var todo Todo
	err = collection.FindOne(context.TODO(), bson.M{"_id": todoID}).Decode(&todo)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	result, err := collection.InsertOne(context.TODO(), todo)
	if err != nil {
		log.Fatal(err)
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	filter := bson.M{"_id": todoID}
	update := bson.M{"$set": bson.M{"task": todo.Task, "completed": todo.Completed}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": todoID}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	collection = client.Database("todo_db").Collection("todos")

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	router.HandleFunc("/login", Login).Methods("POST")

	authenticatedRouter := router.PathPrefix("").Subrouter()
	authenticatedRouter.Use(ValidateToken)

	authenticatedRouter.HandleFunc("/todos", GetTodos).Methods("GET")
	authenticatedRouter.HandleFunc("/todos/{id}", GetTodoByID).Methods("GET")
	authenticatedRouter.HandleFunc("/todos", CreateTodo).Methods("POST")
	authenticatedRouter.HandleFunc("/todos/{id}", UpdateTodoByID).Methods("PUT")
	authenticatedRouter.HandleFunc("/todos/{id}", DeleteTodoByID).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
