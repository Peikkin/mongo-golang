package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Name   string        `json:"name" bson:"_name"`
	Gender string        `json:"gender" bson:"_gender"`
	Age    int           `json:"age" bson:"_age"`
}

type UserController struct {
	Session *mongo.Client
}

func NewUserController(session *mongo.Client) *UserController {
	return &UserController{session}
}

func (client UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	user.ID = bson.NewObjectId()
	collection := client.Session.Database("mongo-golang").Collection("user")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Error().Err(err).Msg("ошибка создания пользователя")
		fmt.Fprintf(w, "ошибка создания пользователя")
		return
	}
	log.Info().Msg("пользователь создан")

	res, err := json.Marshal(user)
	if err != nil {
		log.Error().Err(err).Msg("не удалось прочитать данные")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (client UserController) GetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if !bson.IsObjectIdHex(id) {
		log.Error().Msg("Пользователь не найден")
		fmt.Fprintf(w, "Пользователь не найден")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	objectId := bson.ObjectIdHex(id)

	collection := client.Session.Database("mongo-golang").Collection("user")
	cur, err := collection.Find(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		log.Error().Err(err).Msg("ошибка получения пользователей")
		return
	}
	defer cur.Close(context.Background())

	var users []User
	for cur.Next(context.Background()) {
		var user User
		if err := cur.Decode(&user); err != nil {
			log.Error().Err(err).Msg("ошибка декодирования данных")
			return
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		log.Error().Err(err).Msg("не удалось прочитать данные")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (client UserController) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if !bson.IsObjectIdHex(id) {
		log.Error().Msg("Пользователь не найден")
		fmt.Fprintf(w, "Пользователь не найден")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	objectId := bson.ObjectIdHex(id)

	collection := client.Session.Database("mongo-golang").Collection("user")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		log.Error().Err(err).Msg("ошибка удаления пользователя")
		fmt.Fprintf(w, "ошибка удаления пользователя")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь %v удален!\n", objectId)
}
