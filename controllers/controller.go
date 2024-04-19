package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Peikkin/mongo-golang/models"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	Session *mgo.Session
}

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	objectId := bson.ObjectIdHex(id)
	user := models.User{}
	if err := uc.Session.DB("mongo-golang").C("user").FindId(objectId).One(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось прочитать данные")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	user.ID = bson.NewObjectId()

	if err := uc.Session.DB("mongo-golang").C("user").Insert(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось прочитать данные")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	objectId := bson.ObjectIdHex(id)

	if err := uc.Session.DB("mongo-golang").C("user").RemoveId(objectId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь %v удален!\n", objectId)
}
