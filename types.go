package main

type ProfileResponse struct {
	Name         *string `json:"name"`
	Id           *string `json:"id"`
	Error        *string `json:"error"`
	ErrorMessage *string `json:"errorMessage"`
	Legacy       *bool   `json:"legacy"`
	Demo         *bool   `json:"demo"`
}

type PlayerResponse struct {
	Success bool    `json:"success"`
	Player  *Json   `json:"player"`
	Cause   *string `json:"cause"`
}

type Json map[string]interface{}
