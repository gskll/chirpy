# Chirpy - Golang REST API with http servermux

## Goals

- explore `net/http` servermux
- implement basic refresh token functionality for auth
- use sqlc and goose to handle database querying and migrations

## Features

- we have authenticated users and chirps
- can login using a password and endpoint authentication is using jwt tokens
- with access token user can update user details, create/delete chirps
- access tokens can be refreshed, refresh tokens can be revoked
- have basic 'stripe-like' webhook to upgrade a user to premium status
- get all chirps with user/sorting filters
- page count middleware for 'frontend'

## Quick start

### 1. Setup repo

- clone the repo
- `go mod download`
- `make install-tools`

### 2. Setup environment

- Make sure postgres is installed and running. v14+
- Create a database: `createDb chirpy`

- Clone `.env.dist` to `.env`
  - `DB_URL` is the postgres connection string with ssl disabled: `"postgres://<user>:<pw>@localhost:5432/chirpy?sslmode=disable"`
  - `PLATFORM` should just be `dev`
  - `JWT_SECRET` any random secret to use for jwts
  - `POLKA_KEY` your 'api key' for the polka webhook

NOTE: for the `JWT_SECRET` and `POLKA_KEY` it can be anything. I just generated a random string using `openssl rand -base64 64`

### 3. Run database migrations

- `cd` to `./sql/schema/`
- run `goose POSTGRES_CONNECTION_STRING up`
- you can then check you have `users`, `chirps` and `refresh_tokens` as tables in your db

### 4. Run the server

- `make run` from project root
- server should be running on `localhost:8080`

## Endpoints

! Note: can see examples in the Postman collection file `chirpy.postman_collection.json`

### App

Filserver running on `/app` prefix serving basic html and a logo. Hitting these will increase the page view counter

- `localhost:8080/app`
- `localhost:8080/app/assets/logo/png`

### Admin

#### GET /admin/metrics

Displays HTML page with page view counter for frontend (`/app`)

#### POST /admin/reset

Only available if environment platform is "dev"
Resets the file server hit metrics
Resets the database (deletes everything)

### API

#### GET /api/healthz - Server health check

- Server running
- Response: `200` OK

#### POST /api/users - Create User

- Body: `{email: string, password: string}`
- Response
  - `201`
  - `{
    "id": "ec932e9a-0335-4121-98ab-5ecccb9075d3",
    "created_at": "2024-10-11T15:22:51.955426Z",
    "updated_at": "2024-10-11T15:22:51.955426Z",
    "email": "mike@bettercall.com",
    "is_chirpy_red": false
}`

#### POST /api/login - Login User

- Body:
  `{
  "email": "mike@bettercall.com",
  "password": "98765"
}`
- Response
  - `200`
  - `{
  "id": "ec932e9a-0335-4121-98ab-5ecccb9075d3",
  "created_at": "2024-10-11T15:22:51.955426Z",
  "updated_at": "2024-10-11T15:22:51.955426Z",
  "email": "mike@bettercall.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJlYzkzMmU5YS0wMzM1LTQxMjEtOThhYi01ZWNjY2I5MDc1ZDMiLCJleHAiOjE3Mjg2NTY1NzYsImlhdCI6MTcyODY1Mjk3Nn0.ElLsyeSRHRh7fh_HAjNPC4hpL90WH91qvwrKKkGBdRA",
  "refresh_token": "4cd5437f519a5de25205b47ab3ea0ba70d5365dd93e8b354bbb9b5c7f00fd4f0"
}`

#### POST /api/refresh - Refresh access token

- Auth: Bearer access token
- Response
  - `200`
  - `{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJmZjg4MWQxOS05NmNlLTQ5YTUtODFlYS01MTMxNDc3NzNmYTYiLCJleHAiOjE3Mjg2NTQzODMsImlhdCI6MTcyODY1MDc4M30.PvmLYz1ptimuL9UV8vWQJe0EtE4Hsi-bz-WoL47M5vQ"
}`

#### POST /api/revoke - Revoke refresh token

- Auth: Bearer refresh token
- Response
  - `204`

#### PUT /api/users - Update user email and password

- Auth: Bearer access token
- Body: `{
  "email": "walter@breakingbad.com",
  "password": "j3ssePinkM@nCantCook"
}`
- Response:
  - `200`
  - `{
  "id": "4bb25d3f-0a70-4430-bfb8-2bec1f6c0654",
  "created_at": "2024-10-11T16:46:42.024786Z",
  "updated_at": "2024-10-11T16:46:57.252675Z",
  "email": "walter@breakingbad.com",
  "is_chirpy_red": false
}`

#### POST /api/chirps - Create chirp

- Auth: Bearer access token
- Body: `{
  "body": "Gale!"
}`
- Response:
  - `201`
  - `{
  "id": "f03ea63e-5f29-406f-b19b-a38e127b78bf",
  "user_id": "4a8db05c-e497-4a5e-97b2-7a43a69e2bb5",
  "body": "Gale!",
  "created_at": "2024-10-11T15:23:05.133427Z",
  "updated_at": "2024-10-11T15:23:05.133427Z"
}`

#### GET /api/chirps - Get chirps

- Params:
  - `author_id`: user UUID
  - `sort`: 'asc' or 'desc'
- Example: `GET /api/chirps?sort=asc&author_id=123`
- Response:
  - `200`
  - `[
    {
      "id": "0d634558-e20d-436c-8b18-6b1b9ec15fbc",
      "user_id": "ec932e9a-0335-4121-98ab-5ecccb9075d3",
      "body": "No more half-measures.",
      "created_at": "2024-10-11T15:23:10.007357Z",
      "updated_at": "2024-10-11T15:23:10.007357Z"
    },
    {
      "id": "9b65f428-7661-4b56-8569-ac23594dc1db",
      "user_id": "ec932e9a-0335-4121-98ab-5ecccb9075d3",
      "body": "No more half-measures.",
      "created_at": "2024-10-11T15:23:07.923501Z",
      "updated_at": "2024-10-11T15:23:07.923501Z"
    }
]`

#### GET /api/chirps/{chirpID} - Get chirp

- pathvalue: chirp UUID
- response:
  - `200`
  - `{
  "id": "7c55504d-15ba-4bee-97a7-6793f81b647d",
  "user_id": "4bb25d3f-0a70-4430-bfb8-2bec1f6c0654",
  "body": "Gale!",
  "created_at": "2024-10-11T16:51:46.831441Z",
  "updated_at": "2024-10-11T16:51:46.831441Z"
}`

#### DELETE /api/chirps/{chirpID} - Delete chirp

- Auth: Bearer access token
- Pathvalue: chirp UUID
- Response: `204`

#### POST /api/polka/webhooks - Upgrade user webhook

- Auth: ApiKey polka api key
  - e.g. `Authorization: ApiKey 123`
- Body: `{
  "data": {
    "user_id": "${[ USER_UUID ]}"
  },
  "event": "user.upgraded"
}`
- Response: `204`
