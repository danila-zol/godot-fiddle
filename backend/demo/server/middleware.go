package server

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	operations "game-hangar/database"
)

var client *http.Client

func getClient() *http.Client {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   time.Second,
			ResponseHeaderTimeout: time.Second,
		},
	}
	return &client
}

func SendPostThread(demo *operations.Demo) error {
	var thread operations.Thread
	thread.ID = demo.Thread_id
	thread.Title = demo.Name
	thread.User_id = demo.User_id
	thread.Tag = "demo"
	thread.Topic_id = "topic_7729259c-9059-4ba7-b41b-7ffd4853fa32"
	thread.Created_at, thread.Last_update = demo.Created_at, demo.Updated_at
	thread.Total_upvotes, thread.Total_downvotes = demo.Upvotes, demo.Downvotes
	exportThread, err := json.Marshal(thread)

	if client == nil {
		client = getClient()
	}

	req, err := http.NewRequest("POST", os.Getenv("FORUM_SERVICE_URL"), bytes.NewBuffer(exportThread))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "201" {
		return err
	}
	return nil
}

func SendPatchThread(demo *operations.Demo) error {
	var thread operations.Thread
	thread.ID = demo.Thread_id
	thread.Title = demo.Name
	thread.User_id = demo.User_id
	thread.Tag = "demo"
	thread.Topic_id = "topic_7729259c-9059-4ba7-b41b-7ffd4853fa32"
	thread.Created_at, thread.Last_update = demo.Created_at, demo.Updated_at
	exportThread, err := json.Marshal(thread)
	if err != nil {
		return err
	}

	if client == nil {
		client = getClient()
	}

	req, err := http.NewRequest("PATCH", os.Getenv("FORUM_SERVICE_URL"), bytes.NewBuffer(exportThread))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200" {
		return err
	}
	return nil
}

func SendDeleteThread(id string) error {
	if client == nil {
		client = getClient()
	}

	demo, err := operations.FindFirstDemo(id)
	req, err := http.NewRequest("DELETE", os.Getenv("FORUM_SERVICE_URL")+demo.Thread_id, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200" {
		return err
	}
	return nil
}
