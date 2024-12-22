package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// @Summary	Creates a new user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		User	body		database.User	true	"Create User"
// @Success	200		{object}	ResponseHTTP{data=database.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/users/ [post]
func postUser(w http.ResponseWriter, r *http.Request) {
	var user operations.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostUser handler\n"))
		log.Printf("Error in PostUser handler \n%s", err)
		return
	}

	user.ID = "user_" + uuid.NewString()
	user.Created_at = time.Now()
	user.Karma = 0

	newUser, err := operations.CreateUser(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateUser operation\n"))
		log.Printf("Error in CreateUser operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postUser operation\n"))
		log.Printf("Error in postUser operation \n%s", err)
		return
	}
}

// @Summary	Fetches a user by its ID.
// @Tags		Users
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get User of ID"
// @Success	200	{object}	ResponseHTTP{data=database.User}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/users/{id} [get]
func getUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := operations.FindFirstUser(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: User not found!\n"))
		log.Printf("Error: User not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstUser operation\n"))
		log.Printf("Error in FindFirstUser operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getUserById operation\n"))
		log.Printf("Error in getUserById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all users.
// @Tags		Users
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.User}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/users/ [get]
func getUsers(w http.ResponseWriter, r *http.Request) {
	user, err := operations.FindUsers()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Users not found!\n"))
		log.Printf("Error: Users not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindUsers operation\n"))
		log.Printf("Error in FindUsers operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getUsers operation\n"))
		log.Printf("Error in getUsers operation \n%s", err)
		return
	}
}

// @Summary	Updates a user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		User	body		database.User	true	"Update User"
// @Success	200		{object}	ResponseHTTP{data=database.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/users/ [patch]
func patchUser(w http.ResponseWriter, r *http.Request) {
	var user operations.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchUser handler\n"))
		log.Printf("Error in patchUser handler \n%s", err)
		return
	}

	updUser, err := operations.UpdateUser(user)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Users not found!\n"))
		log.Printf("Error: Users not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateUser operation\n"))
		log.Printf("Error in UpdateUser operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updUser)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchUser handler\n"))
		log.Printf("Error in patchUser handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified user.
// @Tags		Users
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete User of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/users/{id} [delete]
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteUser(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: User not found!\n"))
		log.Printf("Error: User not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstUser operation\n"))
		log.Printf("Error in FindFirstUser operation \n%s", err)
		return
	}
}

// @Summary	Creates a new role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		Role	body		database.Role	true	"Create Role"
// @Success	200		{object}	ResponseHTTP{data=database.Role}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/roles/ [post]
func postRole(w http.ResponseWriter, r *http.Request) {
	var role operations.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostRole handler\n"))
		log.Printf("Error in PostRole handler \n%s", err)
		return
	}

	role.ID = "role_" + uuid.NewString()

	newRole, err := operations.CreateRole(role)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateRole operation\n"))
		log.Printf("Error in CreateRole operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newRole)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postRole operation\n"))
		log.Printf("Error in postRole operation \n%s", err)
		return
	}
}

// @Summary	Fetches a role by its ID.
// @Tags		Roles
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Role of ID"
// @Success	200	{object}	ResponseHTTP{data=database.Role}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/roles/{id} [get]
func getRoleById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	role, err := operations.FindFirstRole(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Role not found!\n"))
		log.Printf("Error: Role not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstRole operation\n"))
		log.Printf("Error in FindFirstRole operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getRoleById operation\n"))
		log.Printf("Error in getRoleById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all roles.
// @Tags		Roles
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.Role}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/roles/ [get]
func getRoles(w http.ResponseWriter, r *http.Request) {
	role, err := operations.FindRoles()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Roles not found!\n"))
		log.Printf("Error: Roles not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindRoles operation\n"))
		log.Printf("Error in FindRoles operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getRoles operation\n"))
		log.Printf("Error in getRoles operation \n%s", err)
		return
	}
}

// @Summary	Updates a role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		Role	body		database.Role	true	"Update Role"
// @Success	200		{object}	ResponseHTTP{data=database.Role}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/roles/ [patch]
func patchRole(w http.ResponseWriter, r *http.Request) {
	var role operations.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchRole handler\n"))
		log.Printf("Error in patchRole handler \n%s", err)
		return
	}

	updRole, err := operations.UpdateRole(role)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Roles not found!\n"))
		log.Printf("Error: Roles not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateRole operation\n"))
		log.Printf("Error in UpdateRole operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updRole)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchRole handler\n"))
		log.Printf("Error in patchRole handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified role.
// @Tags		Roles
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Role of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/roles/{id} [delete]
func deleteRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteRole(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Role not found!\n"))
		log.Printf("Error: Role not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstRole operation\n"))
		log.Printf("Error in FindFirstRole operation \n%s", err)
		return
	}
}
