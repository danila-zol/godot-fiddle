package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"

	_ "game-hangar/docs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// @Summary	Creates a new topic.
// @Tags		Topics
// @Accept		application/json
// @Produce	application/json
// @Param		Topic	body		database.Topic	true	"Create Topic"
// @Success	200		{object}	ResponseHTTP{data=database.Topic}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/topics/ [post]
func postTopic(w http.ResponseWriter, r *http.Request) {
	var topic operations.Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostTopic handler\n"))
		log.Printf("Error in PostTopic handler \n%s", err)
		return
	}

	topic.ID = "topic_" + uuid.NewString()

	newTopic, err := operations.CreateTopic(topic)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateTopic operation\n"))
		log.Printf("Error in CreateTopic operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newTopic)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postTopic operation\n"))
		log.Printf("Error in postTopic operation \n%s", err)
		return
	}
}

// @Summary	Fetches a topic by its ID.
// @Tags		Topics
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Topic of ID"
// @Success	200	{object}	ResponseHTTP{data=database.Topic}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/topics/{id} [get]
func getTopicById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	topic, err := operations.FindFirstTopic(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Topic not found!\n"))
		log.Printf("Error: Topic not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstTopic operation\n"))
		log.Printf("Error in FindFirstTopic operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(topic)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getTopicById operation\n"))
		log.Printf("Error in getTopicById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all topics.
// @Tags		Topics
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.Topic}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/topics/ [get]
func getTopics(w http.ResponseWriter, r *http.Request) {
	topic, err := operations.FindTopics()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Topics not found!\n"))
		log.Printf("Error: Topics not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindTopics operation\n"))
		log.Printf("Error in FindTopics operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(topic)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getTopics operation\n"))
		log.Printf("Error in getTopics operation \n%s", err)
		return
	}
}

// @Summary	Updates a topic.
// @Tags		Topics
// @Accept		application/json
// @Produce	application/json
// @Param		Topic	body		database.Topic	true	"Update Topic"
// @Success	200		{object}	ResponseHTTP{data=database.Topic}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/topics/ [patch]
func patchTopic(w http.ResponseWriter, r *http.Request) {
	var topic operations.Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchTopic handler\n"))
		log.Printf("Error in patchTopic handler \n%s", err)
		return
	}

	updTopic, err := operations.UpdateTopic(topic)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Topics not found!\n"))
		log.Printf("Error: Topics not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateTopic operation\n"))
		log.Printf("Error in UpdateTopic operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updTopic)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchTopic handler\n"))
		log.Printf("Error in patchTopic handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified topic.
// @Tags		Topics
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Topic of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/topics/{id} [delete]
func deleteTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteTopic(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Topic not found!\n"))
		log.Printf("Error: Topic not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstTopic operation\n"))
		log.Printf("Error in FindFirstTopic operation \n%s", err)
		return
	}
}

// @Summary	Creates a new thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		Thread	body		database.Thread	true	"Create Thread"
// @Success	200		{object}	ResponseHTTP{data=database.Thread}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/threads/ [post]
func postThread(w http.ResponseWriter, r *http.Request) {
	var thread operations.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostThread handler\n"))
		log.Printf("Error in PostThread handler \n%s", err)
		return
	}

	thread.ID = "thread_" + uuid.NewString()

	newThread, err := operations.CreateThread(thread)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateThread operation\n"))
		log.Printf("Error in CreateThread operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newThread)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postThread operation\n"))
		log.Printf("Error in postThread operation \n%s", err)
		return
	}
}

// @Summary	Fetches a thread by its ID.
// @Tags		Threads
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Thread of ID"
// @Success	200	{object}	ResponseHTTP{data=database.Thread}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/threads/{id} [get]
func getThreadById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	thread, err := operations.FindFirstThread(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Thread not found!\n"))
		log.Printf("Error: Thread not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstThread operation\n"))
		log.Printf("Error in FindFirstThread operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getThreadById operation\n"))
		log.Printf("Error in getThreadById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all threads.
// @Tags		Threads
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.Thread}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/threads/ [get]
func getThreads(w http.ResponseWriter, r *http.Request) {
	thread, err := operations.FindThreads()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Threads not found!\n"))
		log.Printf("Error: Threads not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindThreads operation\n"))
		log.Printf("Error in FindThreads operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getThreads operation\n"))
		log.Printf("Error in getThreads operation \n%s", err)
		return
	}
}

// @Summary	Updates a thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		Thread	body		database.Thread	true	"Update Thread"
// @Success	200		{object}	ResponseHTTP{data=database.Thread}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/threads/ [patch]
func patchThread(w http.ResponseWriter, r *http.Request) {
	var thread operations.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchThread handler\n"))
		log.Printf("Error in patchThread handler \n%s", err)
		return
	}

	updThread, err := operations.UpdateThread(thread)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Threads not found!\n"))
		log.Printf("Error: Threads not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateThread operation\n"))
		log.Printf("Error in UpdateThread operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updThread)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchThread handler\n"))
		log.Printf("Error in patchThread handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified thread.
// @Tags		Threads
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Thread of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/threads/{id} [delete]
func deleteThread(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteThread(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Thread not found!\n"))
		log.Printf("Error: Thread not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstThread operation\n"))
		log.Printf("Error in FindFirstThread operation \n%s", err)
		return
	}
}

// @Summary	Creates a new message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		Message	body		database.Message	true	"Create Message"
// @Success	200		{object}	ResponseHTTP{data=database.Message}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/messages/ [post]
func postMessage(w http.ResponseWriter, r *http.Request) {
	var message operations.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostMessage handler\n"))
		log.Printf("Error in PostMessage handler \n%s", err)
		return
	}

	message.ID = "message_" + uuid.NewString()

	newMessage, err := operations.CreateMessage(message)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateMessage operation\n"))
		log.Printf("Error in CreateMessage operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newMessage)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postMessage operation\n"))
		log.Printf("Error in postMessage operation \n%s", err)
		return
	}
}

// @Summary	Fetches a message by its ID.
// @Tags		Messages
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Message of ID"
// @Success	200	{object}	ResponseHTTP{data=database.Message}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/messages/{id} [get]
func getMessageById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	message, err := operations.FindFirstMessage(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Message not found!\n"))
		log.Printf("Error: Message not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstMessage operation\n"))
		log.Printf("Error in FindFirstMessage operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getMessageById operation\n"))
		log.Printf("Error in getMessageById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all messages.
// @Tags		Messages
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.Message}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/messages/ [get]
func getMessages(w http.ResponseWriter, r *http.Request) {
	message, err := operations.FindMessages()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Messages not found!\n"))
		log.Printf("Error: Messages not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindMessages operation\n"))
		log.Printf("Error in FindMessages operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getMessages operation\n"))
		log.Printf("Error in getMessages operation \n%s", err)
		return
	}
}

// @Summary	Updates a message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		Message	body		database.Message	true	"Update Message"
// @Success	200		{object}	ResponseHTTP{data=database.Message}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/messages/ [patch]
func patchMessage(w http.ResponseWriter, r *http.Request) {
	var message operations.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchMessage handler\n"))
		log.Printf("Error in patchMessage handler \n%s", err)
		return
	}

	updMessage, err := operations.UpdateMessage(message)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Messages not found!\n"))
		log.Printf("Error: Messages not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateMessage operation\n"))
		log.Printf("Error in UpdateMessage operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updMessage)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchMessage handler\n"))
		log.Printf("Error in patchMessage handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified message.
// @Tags		Messages
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Message of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/messages/{id} [delete]
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteMessage(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Message not found!\n"))
		log.Printf("Error: Message not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstMessage operation\n"))
		log.Printf("Error in FindFirstMessage operation \n%s", err)
		return
	}
}
