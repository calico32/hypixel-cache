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

func fetchFromApi(c *gin.Context) {
	if c.Writer.Written() {
		c.Next()
		return
	}

	identifierType := c.Param("type")
	identifier := c.Param("identifier")

	var uuid string

	if identifierType != "uuid" && identifierType != "name" {
		finish(c, 400, NewFailure("invalid identifier type"))
		return
	}

	if identifierType == "uuid" {
		if !uuidRegexp.MatchString(identifier) {
			finish(c, 400, NewFailure("invalid uuid"))
			return
		}

		uuid = identifier
	} else if identifierType == "name" {
		if !usernameRegexp.MatchString(identifier) {
			finish(c, 400, NewFailure("invalid username"))
			return
		}

		if cachedUuid, ok := uuidCache.Get(identifier); ok {
			uuid = cachedUuid.(string)
		} else {
			profile, err := fetchProfile(strings.ToLower(identifier))
			if err != nil {
				if err.Error() == "user not found" {
					finish(c, 200, NewSuccessNotFound(time.Now(), false))
					return
				}

				finish(c, 500, NewFailure("error fetching profile: "+err.Error()))
				return
			}

			uuid = *profile.Id
		}
	}

	uuid = removeDashes(uuid)

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

func fetchProfile(username string) (profileResponse *ProfileResponse, err error) {
	if !usernameRegexp.MatchString(username) {
		err = errors.New("invalid username")
		return
	}

	res, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)

	if err != nil {
		return
	}

	if res.StatusCode == 404 || res.StatusCode == 204 {
		err = errors.New("user not found")
		return
	} else if res.StatusCode != 200 {
		err = errors.New("error fetching profile: " + res.Status)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &profileResponse)
	if err != nil {
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
		return
	}

	if res.StatusCode == 404 || profileResponse.Demo != nil || profileResponse.Id == nil {
		return
	}

	if !strings.EqualFold(username, *profileResponse.Name) {
		err = errors.New("provided and resolved username mismatch: " + username + " != " + *profileResponse.Name)
		return
	}

	uuidCache.Set(strings.ToLower(username), removeDashes(*profileResponse.Id), cache.DefaultExpiration)

	return
}
