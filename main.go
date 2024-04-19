package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Peikkin/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Устанавливаем параметры подключения к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Подключаемся к MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось подключиться к MongoDB")
	}

	// Проверяем подключение к MongoDB
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подключения к MongoDB")
	}

	log.Info().Msg("Подключение к MongoDB успешно!")

	// Закрываем соединение с MongoDB
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Ошибка закрытия подключения к MongoDB")
		}
	}()

	router := httprouter.New()
	session := controllers.NewUserController(client)
	router.GET("/user/:id", session.GetUsers)
	router.POST("/user", session.CreateUser)
	router.DELETE("/user/:id", session.DeleteUser)

	log.Info().Msg("Запуск сервера на порту :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("Ошибка запуска сервера")
	}
}

// func getSession() *mongo.Client {
// 	session, err := mgo.Dial("mongodb://localhost:27017")
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Не удалось подключиться к MongoDB")
// 	}
// 	log.Info().Msg("Подключение к MongoDB успешно")
// 	return session.Copy()
// }
