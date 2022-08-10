package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jumaniyozov/gores/internal/app/models"
	"net/http"
	"strconv"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func MsgEncoder(w http.ResponseWriter, value interface{}, api *API, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		api.logger.Fatal("Error while encoding message")
	}
}

func (api *API) GetAllArticles(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	articles, err := api.storage.Article().SelectAll()
	if err != nil {
		api.logger.Info("Error occurred while Articles.SelectAll:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "Internal server error occurred. Try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	api.logger.Info("Get All Articles GET /api/v1/articles")
	MsgEncoder(writer, articles, api, 200)
}

func (api *API) GetArticleByID(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Article by ID /api/v1/articles/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unapropriate id value. Don't use ID as uncasting to int value",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	article, ok, err := api.storage.Article().FindArticleByID(id)
	if err != nil {
		api.logger.Info("Trouble while accessing DB table(articles) with id. Err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Internal server error occurred. Please, try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Article with that ID does not exist",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}
	MsgEncoder(writer, article, api, 200)
}

func (api *API) DeleteArticleByID(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Article by ID /api/v1/articles/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unapropriate id value. Don't use ID as uncasting to int value",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	_, ok, err := api.storage.Article().FindArticleByID(id)
	if err != nil {
		api.logger.Info("Trouble while accessing DB table(articles) with id. Err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Internal server error occurred. Please, try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Article with that ID does not exist",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	_, err = api.storage.Article().DeleteByID(id)
	if err != nil {
		api.logger.Info("Trouble while deleting article from DB table(articles) with id. Err:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "Internal server error occurred. Please, try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	msg := Message{
		StatusCode: 202,
		Message:    fmt.Sprintf("Article with ID %d successfully deleted", id),
		IsError:    false,
	}

	MsgEncoder(writer, msg, api, msg.StatusCode)
}

func (api *API) PostArticle(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Article POST /api/v1/articles")

	var article models.Article
	if err := json.NewDecoder(req.Body).Decode(&article); err != nil {
		api.logger.Info("Invalid json received from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, 400)
		return
	}

	a, err := api.storage.Article().Create(&article)
	if err != nil {
		api.logger.Info("Trouble while creating new article, ", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troulbe accessing database. Try again later. ",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, 501)
		return
	}
	MsgEncoder(writer, a, api, 201)
}

func (api *API) PostUserRegister(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /api/v1/user/register")

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		api.logger.Info("Trouble while creating new user.", err)
		msg := Message{
			StatusCode: 404,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	_, ok, err := api.storage.User().FindByLogin(user.Login)
	if err != nil {
		api.logger.Info("Trouble while accessing DB table(users) with id. Err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Internal server error occurred. Please, try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	if ok {
		api.logger.Info("User with that login already exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User with that login already exists",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}
	userAdded, err := api.storage.User().Create(&user)
	if err != nil {
		api.logger.Info("Trouble while creating user with credentials. Err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Internal server error occurred. Please, try again later",
			IsError:    true,
		}
		MsgEncoder(writer, msg, api, msg.StatusCode)
		return
	}

	msg := Message{
		StatusCode: 201,
		Message:    fmt.Sprintf("User {login:%s} successfully registered", userAdded.Login),
		IsError:    false,
	}
	MsgEncoder(writer, msg, api, msg.StatusCode)
}
