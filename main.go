package main

import (
	"net/http"
	"os"

	"github.com/Peikkin/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	router := httprouter.New()
	session := controllers.NewUserController(getSession())
	router.GET("/user/:id", session.GetUser)
	router.POST("/user", session.CreateUser)
	router.DELETE("/user/:id", session.DeleteUser)

	log.Info().Msg("Запуск сервера на порту :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("Ошибка запуска сервера")
	}
}

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось подключиться к MongoDB")
	}
	log.Info().Msg("Подключение к MongoDB успешно")
	return session.Copy()
}
