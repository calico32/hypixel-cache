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
	InvalidIdentType = "invalid type"
	InvalidUuid      = "invalid uuid"
	InvalidName      = "invalid username"
	ProfileNotFound  = "profile not found"
	ServerError      = "server error"
)

func getUuid(typ string, ident string) (uuid string, err error, errorCode int) {
	ident = strings.ToLower(ident)

	if typ != "uuid" && typ != "name" {
		err = errors.New(InvalidIdentType)
		errorCode = 400
		return
	}

	if typ == "uuid" {
		if !uuidRegexp.MatchString(ident) {
			err = errors.New(InvalidUuid)
			errorCode = 400
			return
		}

		uuid = ident
	} else if typ == "name" {
		if !usernameRegexp.MatchString(ident) {
			err = errors.New(InvalidName)
			errorCode = 400
			return
		}

		if cachedUuid, ok := uuidCache.Get(ident); ok {
			uuid = cachedUuid.(string)
		} else {
			profile, profileErr, profileErrorCode := fetchProfile(strings.ToLower(ident))
			if profileErr != nil {
				err = profileErr
				errorCode = profileErrorCode
				return
			}

			uuid = *profile.Id
		}
	}

	uuid = removeDashes(uuid)
	return
}
