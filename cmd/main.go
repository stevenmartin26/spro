package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	db, err := newDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := repository.NewUserRepository(repository.UserRepositoryImplOptions{
		DB: db,
	})
	loginLogRepository := repository.NewLoginLogRepositoryImpl(repository.LoginLogRepositoryImplOptions{
		DB: db,
	})

	tokenManager := service.NewJWTManager(os.Getenv("PUBLIC_KEY_PATH"), os.Getenv("PRIVATE_KEY_PATH"))
	authService := service.NewAuthServiceImpl(userRepository, loginLogRepository, tokenManager)
	profileService := service.NewProfileServiceImpl(userRepository, tokenManager)

	opts := handler.NewServerOptions{
		AuthService:    authService,
		ProfileService: profileService,
	}
	return handler.NewServer(opts)
}

func newDatabase() (*sql.DB, error) {
	dbDsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		log.Fatal("error opening database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("error pinging database:", err)
	}
	fmt.Println("connected to the database!")

	return db, nil
}
