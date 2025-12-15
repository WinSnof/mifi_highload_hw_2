package handlers

import (
	"homework_2/models"
	"homework_2/services"
	"homework_2/utils" // Используем наш универсальный пакет
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var UserService = services.NewUserService()

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := UserService.GetAll()
	utils.RespondWithJSON(w, http.StatusOK, users)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, ok := UserService.GetByID(id)
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// CreateUserHandler обрабатывает POST /api/users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// --- УНИВЕРСАЛЬНАЯ ПРОВЕРКА JSON И ВАЛИДАЦИЯ ---
	// &user реализует utils.Validatable
	if utils.DecodeAndValidate(w, r, &user) {
		return // Ошибка (JSON или валидация) уже обработана и отправлена 400
	}
	// ------------------------------------------------

	// Сохранение пользователя
	savedUser := UserService.Create(user)

	// Асинхронное логирование
	go utils.LogUserAction("CREATE", savedUser.ID)

	utils.RespondWithJSON(w, http.StatusCreated, savedUser)
}

// UpdateUserHandler обрабатывает PUT /api/users/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User

	// --- УНИВЕРСАЛЬНАЯ ПРОВЕРКА JSON И ВАЛИДАЦИЯ ---
	if utils.DecodeAndValidate(w, r, &user) {
		return // Ошибка (JSON или валидация) уже обработана и отправлена 400
	}
	// ------------------------------------------------

	updatedUser, ok := UserService.Update(id, user)
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Асинхронное логирование (goroutine)
	go utils.LogUserAction("UPDATE", id)

	utils.RespondWithJSON(w, http.StatusOK, updatedUser)
}

// DeleteUserHandler обрабатывает DELETE /api/users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if !UserService.Delete(id) {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Асинхронное логирование (goroutine)
	go utils.LogUserAction("DELETE", id)

	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
