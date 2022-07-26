package main

import (
	"errors"
	"strings"
)

func removeDashes(s string) string {
	return strings.Replace(s, "-", "", -1)
}

func insertUuidDashes(s string) string {
	s = removeDashes(s)

	if len(s) != 32 {
		panic("Invalid UUID")
	}

	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

const (
	InvalidIdentType = "invalid identifier type"
	InvalidUuid      = "invalid uuid"
	InvalidName      = "invalid username"
	ProfileNotFound  = "profile not found"
	PlayerNotFound   = "player not found"
	ServerError      = "server error"
)

func getUuid(typ string, ident string) (uuid string, err error) {

	if typ != "uuid" && typ != "name" {
		err = errors.New(InvalidIdentType)
		return
	}

	if typ == "uuid" {
		if !uuidRegexp.MatchString(ident) {
			err = errors.New(InvalidUuid)
			return
		}

		uuid = ident
	} else if typ == "name" {
		if !usernameRegexp.MatchString(ident) {
			err = errors.New(InvalidName)
			return
		}

		if cachedUuid, ok := uuidCache.Get(ident); ok {
			uuid = cachedUuid.(string)
		} else {
			profile, profileErr := fetchProfile(strings.ToLower(ident))
			if profileErr != nil {
				err = profileErr
				return
			}

			uuid = *profile.Id
		}
	}

	uuid = removeDashes(uuid)
	return
}
