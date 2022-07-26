# hypixel-cache

A simple cache/intermediary for the [Hypixel API](https://api.hypixel.net/) written in Go.

## Configuration

See [.env.example](./.env.example). Copy to `.env` to use.

## API

`hypixel-cache` supports fetching player data by both UUIDs (which are sent directly to the Hypixel API) and usernames (which are resolved to UUIDs first via a Mojang API).

Hypixel player data is kept around for 15-30 minutes. Subsequent requests for the same player will return the cached data. 

Similarly, UUIDs that have been resolved from usernames will be cached for 60-120 minutes. Subsequent requests for the same username will bypass the Mojang API and use the cached UUID. This means if a player changes their username during those 60-120 minutes, both the new and old username will resolve to their UUID.

### Definitions

**API key**: This is a separate, user-defined key for the cache, a.k.a. **NOT** the Hypixel API key. The Hypixel API key is specified in the environment variable `HYPIXEL_API_KEY`.

**Account/Profile**: Minecraft account. If the account exists, then there is someone with the given username/UUID.

**Player**: Hypixel player data. If the player exists, then the account exists and has logged in to Hypixel at least once.

### Requests
```http
GET /uuid/b41d72aa-0098-45b0-ada5-f154d8796d20
X-Secret: {api_key}
```
```http
GET /uuid/b41d72aa009845b0ada5f154d8796d20
X-Secret: {api_key}
```
```http
GET /name/wiisportsresorts
X-Secret: {api_key}
```

### Responses

#### Account exists, player exists
```http
HTTP/1.1 200 OK
Access-Control-Allow-Methods: GET
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: X-Response-Time
X-Response-Time: 123ms
Content-Type: application/json

{
  "success": true,
  "fetchedAt": "1970-01-01T00:00:00Z",
  "cached": false,
  "username": "…", // mirrors `player.displayname` (empty string if not present)
  "uuid": "…",     // mirrors `player.uuid`
  "player": {
    // data as returned by the Hypixel API
  }
}
```


#### Account does not exist (`/name/…` only)
```http
HTTP/1.1 404 Not Found
...
Content-Type: application/json

{
  "success": false,
  "error": "profile not found"
}
```

#### Account does not exist (`/uuid/…` only); Account exists, player does not exist (either endpoint)
```http
HTTP/1.1 200 OK
...
Content-Type: application/json

{
  "success": true,
  "fetchedAt": "1970-01-01T00:00:00Z",
  "cached": false
}
```

#### Invalid request
```http
HTTP/1.1 400 Bad Request
...
Content-Type: application/json

{
  "success": false,
  "error": "…",
  // "invalid type"     request path was not /name/… or /uuid/…
  // "invalid uuid"     must match /^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$/i
  // "invalid name"     must match /^[a-zA-Z0-9_]{3,16}$/
}
```

#### Invalid or missing API key
```http
HTTP/1.1 401 Unauthorized
...
Content-Type: application/json

{
  "success": false,
  "error": "unauthorized"
}
```

#### Hypixel API ratelimit
```http
HTTP/1.1 429 Too Many Requests
...
Content-Type: application/json

{
  "success": false,
  "error": "ratelimited, try again later"
}
```

#### Server error
```http
HTTP/1.1 500 Internal Server Error
...
Content-Type: application/json

{
  "success": false,
  "error": "…"  // varies
}
```
