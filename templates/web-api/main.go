package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"PROJECT_NAME/async"
	"PROJECT_NAME/types"
	"PROJECT_NAME/utils"

	"github.com/gorilla/mux"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// In-memory user store (use database in production)
var users = types.NewMap[int, *User]()
var userIdCounter = 1

// Event emitter for user events
var userEvents = types.NewEventEmitter[UserEvent]()

type UserEvent struct {
	Type string      `json:"type"`
	User *User       `json:"user"`
	Time time.Time   `json:"time"`
}

func main() {
	fmt.Println("🚀 Starting TypeScript-like Go Web API...")
	
	// Initialize some sample data
	initializeData()
	
	// Setup event listeners
	setupEventListeners()
	
	// Setup routes
	router := setupRoutes()
	
	// Start server
	fmt.Println("🌐 Server starting on http://localhost:8080")
	fmt.Println("📡 Available endpoints:")
	fmt.Println("  GET    /users       - List all users")
	fmt.Println("  GET    /users/{id}  - Get user by ID")
	fmt.Println("  POST   /users       - Create new user")
	fmt.Println("  PUT    /users/{id}  - Update user")
	fmt.Println("  DELETE /users/{id}  - Delete user")
	fmt.Println("  GET    /health      - Health check")
	fmt.Println()
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeData() {
	// Add some sample users
	users.Set(1, &User{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", Age: 28})
	users.Set(2, &User{ID: 2, Name: "Bob Smith", Email: "bob@example.com", Age: 32})
	users.Set(3, &User{ID: 3, Name: "Carol Brown", Email: "carol@example.com", Age: 25})
	userIdCounter = 4
}

func setupEventListeners() {
	userEvents.On("user:created", func(event UserEvent) {
		fmt.Printf("📢 Event: User %s was created\n", event.User.Name)
	})
	
	userEvents.On("user:updated", func(event UserEvent) {
		fmt.Printf("📢 Event: User %s was updated\n", event.User.Name)
	})
	
	userEvents.On("user:deleted", func(event UserEvent) {
		fmt.Printf("📢 Event: User %s was deleted\n", event.User.Name)
	})
}

func setupRoutes() *mux.Router {
	router := mux.NewRouter()
	
	// Middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	
	// Routes
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/users", getUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	router.HandleFunc("/users", createUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")
	
	return router
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("📝 %s %s - %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Handlers
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"users":     users.Size(),
	}
	
	writeJSONResponse(w, http.StatusOK, response)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Use async operation to simulate database query
	promise := async.NewPromise(func() ([]User, error) {
		time.Sleep(10 * time.Millisecond) // Simulate DB query
		
		allUsers := make([]User, 0, users.Size())
		users.ForEach(func(id int, user *User) {
			allUsers = append(allUsers, *user)
		})
		
		return allUsers, nil
	})
	
	result, err := promise.Await()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}
	
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"users": result,
		"total": len(result),
	})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	
	userOpt := users.Get(id)
	if userOpt.IsNone() {
		writeErrorResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}
	
	writeJSONResponse(w, http.StatusOK, userOpt.Get())
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}
	
	// Validate user
	if err := validateUser(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}
	
	// Assign ID and save
	user.ID = userIdCounter
	userIdCounter++
	users.Set(user.ID, &user)
	
	// Emit event
	userEvents.Emit("user:created", UserEvent{
		Type: "created",
		User: &user,
		Time: time.Now(),
	})
	
	writeJSONResponse(w, http.StatusCreated, user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	
	existingUserOpt := users.Get(id)
	if existingUserOpt.IsNone() {
		writeErrorResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}
	
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}
	
	// Validate user
	if err := validateUser(&updatedUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}
	
	updatedUser.ID = id
	users.Set(id, &updatedUser)
	
	// Emit event
	userEvents.Emit("user:updated", UserEvent{
		Type: "updated",
		User: &updatedUser,
		Time: time.Now(),
	})
	
	writeJSONResponse(w, http.StatusOK, updatedUser)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	
	userOpt := users.Get(id)
	if userOpt.IsNone() {
		writeErrorResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}
	
	user := userOpt.Get()
	users.Delete(id)
	
	// Emit event
	userEvents.Emit("user:deleted", UserEvent{
		Type: "deleted",
		User: user,
		Time: time.Now(),
	})
	
	writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

// Utility functions
func validateUser(user *User) error {
	// Use TypeScript-like validation
	return types.NewTry(func() (interface{}, error) {
		if utils.Strings.IsBlank(user.Name) {
			return nil, types.NewValidationError("Name is required")
		}
		if utils.Strings.IsBlank(user.Email) {
			return nil, types.NewValidationError("Email is required")
		}
		if user.Age < 0 || user.Age > 150 {
			return nil, types.NewValidationError("Age must be between 0 and 150")
		}
		if !utils.Strings.Contains(user.Email, "@") {
			return nil, types.NewValidationError("Invalid email format")
		}
		return nil, nil
	}).Execute()
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("❌ Error encoding JSON: %v\n", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	errorResponse := map[string]interface{}{
		"error":     message,
		"timestamp": time.Now(),
	}
	
	if err != nil {
		errorResponse["details"] = err.Error()
	}
	
	writeJSONResponse(w, status, errorResponse)
}