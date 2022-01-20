package main

import (
	"Technopark_DB_Project/app/handlers"
	"Technopark_DB_Project/app/repositories/stores"
	"Technopark_DB_Project/app/usecases/impl"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	_ "github.com/lib/pq"
)

type Server struct {
	settings Settings
}

func CreateServer() *Server {
	settings := InitSettings()
	return &Server{settings: settings}
}

func (server *Server) Run() {
	router := gin.New()

	// Postgres
	conn, err := pgx.ParseConnectionString(server.settings.PostgresDsn)
	if err != nil {
		fmt.Println(err)
	}
	postgresConnection, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 1000,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer postgresConnection.Close()

	// Repositories
	userRepo := stores.CreateUserRepository(postgresConnection)
	forumRepo := stores.CreateForumRepository(postgresConnection)
	postRepo := stores.CreatePostRepository(postgresConnection)
	serviceRepo := stores.CreateServiceRepository(postgresConnection)
	threadRepo := stores.CreateThreadRepository(postgresConnection)
	voteRepo := stores.CreateVoteRepository(postgresConnection)

	// UseCases
	userUseCase := impl.CreateUserUseCase(userRepo)
	forumUseCase := impl.CreateForumUseCase(forumRepo, threadRepo, userRepo)
	postUseCase := impl.CreatePostUseCase(postRepo, userRepo, threadRepo, forumRepo)
	serviceUseCase := impl.CreateServiceUseCase(serviceRepo)
	threadUseCase := impl.CreateThreadUseCase(threadRepo, voteRepo)

	// Middlewares
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.CorsConfig))

	// Handlers
	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateUserHandler(rootGroup, server.settings.UserURL, userUseCase)
	handlers.CreateForumHandler(rootGroup, server.settings.ForumURL, forumUseCase)
	handlers.CreatePostHandler(rootGroup, server.settings.PostURL, postUseCase)
	handlers.CreateServiceHandler(rootGroup, server.settings.ServiceURL, serviceUseCase)
	handlers.CreateThreadHandler(rootGroup, server.settings.ThreadURL, threadUseCase)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		fmt.Println(err)
	}
}
