package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var apiResolvedUuids = make(map[*http.Request]string)

func fetchFromApi(c *gin.Context) {
	if c.Writer.Written() {
		c.Next()
		return
	}

	uuid, ok := apiResolvedUuids[c.Request]
	delete(apiResolvedUuids, c.Request)
	if !ok {
		newUuid, err, code := getUuid(c.Param("type"), c.Param("identifier"))
		if err != nil {
			finish(c, code, NewFailure(err.Error()))
			return
		}
		uuid = newUuid
	}

	requestUrl := "https://api.hypixel.net/player?key=" + os.Getenv("HYPIXEL_API_KEY") + "&uuid=" + uuid

	res, err := http.Get(requestUrl)
	if err != nil {
		finish(c, 500, NewFailure("error fetching player: "+err.Error()))
		return
	}

	if res.StatusCode == 429 {
		finish(c, 429, NewFailure("ratelimited, try again later"))
		return
	}

	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		finish(c, 500, NewFailure("error parsing response: "+err.Error()))
		return
	}

	if response["player"] == nil {

		cached := CachedPlayer{
			FetchedAt: time.Now(),
		}
		playerCache.Set(uuid, cached, cache.DefaultExpiration)
		finish(c, 200, NewSuccessNotFound(cached.FetchedAt, false))
		return
	}

	player := response["player"].(map[string]interface{})

	cached := CachedPlayer{
		Player:    player,
		FetchedAt: time.Now(),
	}
	playerCache.Set(uuid, cached, cache.DefaultExpiration)
	finish(c, 200, NewSuccessPlayerFound(cached, false))
}

func fetchProfile(username string) (profileResponse *ProfileResponse, err error, errorCode int) {
	if !usernameRegexp.MatchString(username) {
		err = errors.New(string(InvalidName))
		errorCode = 400
		return
	}

	res, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)

	if err != nil {
		return
	}

	if res.StatusCode == 404 || res.StatusCode == 204 {
		err = errors.New(string(ProfileNotFound))
		errorCode = 404
		return
	} else if res.StatusCode != 200 {
		err = errors.New("error fetching profile: " + res.Status)
		errorCode = 500
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errorCode = 500
		return
	}

	err = json.Unmarshal(body, &profileResponse)
	if err != nil {
		errorCode = 500
		return
	}

	if profileResponse.Error != nil {
		message := *profileResponse.Error + ": "
		if profileResponse.ErrorMessage != nil {
			message += *profileResponse.ErrorMessage
		} else {
			message += "(unknown)"
		}

		err = errors.New(message)
		errorCode = 500
		return
	}

	if res.StatusCode == 404 || profileResponse.Demo != nil || profileResponse.Id == nil {
		return
	}

	if !strings.EqualFold(username, *profileResponse.Name) {
		err = errors.New("provided and resolved username mismatch: " + username + " != " + *profileResponse.Name)
		errorCode = 500
		return
	}

	uuidCache.Set(strings.ToLower(username), removeDashes(*profileResponse.Id), cache.DefaultExpiration)
	return
}
