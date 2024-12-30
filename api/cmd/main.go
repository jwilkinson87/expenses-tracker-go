package cmd

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"example.com/expenses-tracker/api/internal/handlers"
	"example.com/expenses-tracker/api/internal/http"
	"example.com/expenses-tracker/api/internal/http/middleware"
	"example.com/expenses-tracker/api/internal/repositories"
	"example.com/expenses-tracker/api/internal/validation"
	"example.com/expenses-tracker/pkg/database"
	"example.com/expenses-tracker/pkg/encryption"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Container struct {
	UserRepository     repositories.UserRepository
	UserAuthRepository repositories.UserAuthRepository
	ExpenseRepository  repositories.ExpenseRepository
	AuthHandler        *handlers.AuthHandler
	EncryptionHandler  *encryption.EncryptionHandler
	middleware         map[string]gin.HandlerFunc
}

var (
	container *Container
	once      sync.Once
)

const (
	errFailedToConnectToDatabase = "failed to connect to database: %w"
)

// Setup prepares this application
func Setup() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	engine := gin.New()
	engine.Use(gin.Recovery())

	db, err := database.NewDatabase()
	if err != nil {
		panic(fmt.Errorf(errFailedToConnectToDatabase, err))
	}

	container = &Container{}
	once.Do(func() {
		setupContainer(db)
	})

	setupMiddleware(engine)
	setupHttpHandlers(engine)
	setupValidators()

	engine.Run()
	defer db.Close()
}

func setupContainer(db *sql.DB) {
	container.UserAuthRepository = repositories.NewAuthRepository(db)
	container.UserRepository = repositories.NewUserRepository(db)
	container.ExpenseRepository = repositories.NewExpensesRepository(db)
	container.EncryptionHandler = encryption.NewEncryptionHandlerFromEnvVars()
	container.AuthHandler = handlers.NewAuthHandler(container.UserAuthRepository, container.UserRepository, container.EncryptionHandler)
	container.middleware = make(map[string]gin.HandlerFunc, 1)
}

func setupMiddleware(g *gin.Engine) {
	g.Use(middleware.RequestIdMiddleware())
	g.Use(middleware.LoggerMiddleware())

	authMiddleware := middleware.NewAuthMiddleware(container.AuthHandler)
	container.middleware["auth"] = authMiddleware.HandleAuthToken()
}

func setupHttpHandlers(g *gin.Engine) {
	expensesGroup := g.Group("/api/expenses")
	expensesGroup.Use(container.middleware["auth"])
	expenseHandler := http.NewExpensesHandler(container.ExpenseRepository)
	expenseHandler.RegisterRoutes(expensesGroup)

	userHandler := http.NewUsersHandler(container.UserRepository)
	userHandler.RegisterRoutes(g)

	authHandler := http.NewAuthHandler(*container.AuthHandler)
	authHandler.RegisterRoutes(g)
}

func setupValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validpassword", validation.ValidPassword)
	}
}
