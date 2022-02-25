package main

import "time"

type successPlayerFound struct {
	Success   bool   `json:"success"`
	Player    Json   `json:"player"`
	FetchedAt string `json:"fetchedAt"`
	Cached    bool   `json:"cached"`
	Username  string `json:"username"`
	Uuid      string `json:"uuid"`
}

func NewSuccessPlayerFound(data CachedPlayer, cached bool) successPlayerFound {
	return successPlayerFound{
		Success:   true,
		Player:    data.Player,
		FetchedAt: data.FetchedAt.Format(time.RFC3339),
		Cached:    cached,
		Username:  data.Player["displayname"].(string),
		Uuid:      data.Player["uuid"].(string),
	}
}

type successNotFound struct {
	Success   bool   `json:"success"`
	FetchedAt string `json:"fetchedAt"`
	Cached    bool   `json:"cached"`
}

func NewSuccessNotFound(fetchedAt time.Time, cached bool) successNotFound {
	return successNotFound{
		Success:   true,
		FetchedAt: fetchedAt.Format(time.RFC3339),
		Cached:    cached,
	}
}

type failure struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewFailure(reason string) failure {
	return failure{
		Success: false,
		Error:   reason,
	}
}
