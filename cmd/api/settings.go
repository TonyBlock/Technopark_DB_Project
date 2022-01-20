package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
)

type Settings struct {
	RootURL    string
	ForumURL   string
	PostURL    string
	ThreadURL  string
	UserURL    string
	ServiceURL string

	ServerAddress string

	Origins        []string
	AllowedMethods []string

	CorsConfig cors.Config

	dbPort     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbName     string

	PostgresDsn string
}

func InitSettings() (settings Settings) {
	settings = Settings{
		RootURL:    "/api",
		ForumURL:   "/forum",
		PostURL:    "/post",
		ThreadURL:  "/thread",
		UserURL:    "/user",
		ServiceURL: "/service",

		ServerAddress: ":5000",

		Origins: []string{
			"http://localhost:5000",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		dbPort:     "5432",
		dbUser:     "Forum_user",
		dbPassword: "db_password",
		dbHost:     "localhost",
		dbName:     "db_forum",

		CorsConfig: cors.DefaultConfig(),
	}

	settings.PostgresDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", settings.dbHost, settings.dbUser, settings.dbPassword, settings.dbName, settings.dbPort)

	settings.CorsConfig.AllowOrigins = settings.Origins
	settings.CorsConfig.AllowCredentials = true

	return
}
