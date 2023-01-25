# telegram-logger

[![codecov](https://codecov.io/gh/psyb0t/telegram-logger/branch/master/graph/badge.svg?token=FBYNVHPF8Q)](https://codecov.io/gh/psyb0t/telegram-logger)
[![goreportcard](https://goreportcard.com/badge/github.com/psyb0t/telegram-logger)](https://goreportcard.com/report/github.com/psyb0t/telegram-logger)
[![pipeline](https://github.com/psyb0t/telegram-logger/actions/workflows/pipeline.yml/badge.svg)](https://github.com/psyb0t/telegram-logger/actions/workflows/pipeline.yml)

This service is a server that receives log entries via HTTP POST requests and sends them to a Telegram bot. The log entries are expected to be in the form of a JSON object with various fields representing log data and metadata.

## Configuration

The configuration file is a YAML file that specifies the following fields:

- ListenAddress: the address that the HTTP server should listen on (e.g. "localhost:8080"). Defaults to "0.0.0.0:80".
- Logger: a nested object for configuring the logger.
  - Level: the log level (e.g. "debug", "info", "warning", "error"). Defaults to "debug".
  - Format: the log format (e.g. "text", "json"). Defaults to "json".
- TelegramBot: a nested object for configuring the Telegram bot.
  - Token: the token used to authenticate the connection to the Telegram bot.
  - SuperuserChatID: the chat ID of the superadmin user.
- Storage: a nested object for configuring the storage backend.
  - Type: the type of storage backend to use (e.g. "badgerDB").
  - BadgerDB: a nested object for configuring the BadgerDB storage backend.
    - DSN: the Data Source Name (DSN) used to connect to the BadgerDB instance (e.g. "mydatabase").

The configuration file can be specified using the `CONFIGFILE` environment variable or, if not set, the default file `./config.yml` will be used. If the default file does not exist it will try working with the given environment variables (see bellow).

Here's an example configuration file:

```yaml
listenAddress: 0.0.0.0:8080
logger:
  level: debug
  format: json
telegramBot:
  token:
  superuserChatID: 38081130
storage:
  type: badgerDB
  badgerDB:
    dsn: /path/to/db/dir
```

All of the settings from the configuration file can also be set via environment variables. The environment variables override any value already present in the config file.

- `LISTENADDRESS`: overrides the listenAddress field in the config struct.
- `LOGGER_LEVEL`: overrides the level field in the logger struct.
- `LOGGER_FORMAT`: overrides the format field in the logger struct.
- `TELEGRAMBOT_TOKEN`: overrides the token field in the telegramBot struct.
- `TELEGRAMBOT_SUPERUSERCHATID`: overrides the superuserChatID field in the telegramBot struct.
- `STORAGE_TYPE`: overrides the type field in the storage struct.
- `STORAGE_BADGERDB_DSN`: overrides the dsn field in the badgerDB struct.

## HTTP API

The service listens for HTTP POST requests to the root path and sends the request body as a message via the Telegram bot to the user associated with the request (determined by the value of the `X-ID` header). The response body is a JSON object.

### Request

```json
{
  "caller": "myService",
  "time": "2022-03-11T12:34:56.789Z",
  "level": "error",
  "message": "Something went wrong!",
  "error": "Error: XYZ",
  "requestID": "abc123",
  "traceID": "def456",
  "spanID": "ghi789",
  "data": {
    "someKey": "someValue",
    "anotherKey": 123
  }
}
```

### Success response

```json
{
  "message": "successfully sent log entry via Telegram"
}
```

### Error response

```json
{
  "error": "Bad Request: message text is empty"
}
```

## Telegram Bot

The service has a number of commands that can be sent via the Telegram bot by a user. These commands are handled by the app and invoke the corresponding command handler function. The following commands are supported:

- /start: generates a unique ID for the user and sends it to the user in a welcome message. The user can then use this ID in the X-ID header of HTTP POST requests to the app to authenticate themselves.
- /stop: stops sending log entries to the user and removes the user's ID from the database.
- /getAllUsers (superadmin only): retrieves a list of all users from the database and sends it to the user via the Telegram bot.

## Running the service

When running via using the docker image you will need to use port 80 on the listen address because that's the only one being exposed so you can just not specify it anywhere.

### Run via docker

```bash
docker run -it \
    --env CONFIGFILE=/app/config.yml \
    --env LOGGER_LEVEL=debug \
    --env LOGGER_FORMAT=json \
    --env TELEGRAMBOT_TOKEN=TOKEN \
    --env TELEGRAMBOT_SUPERUSERCHATID=CHATID \
    --env STORAGE_TYPE=badgerDB \
    --env STORAGE_BADGERDB_DSN=/app/db \
    --mount type=bind,src=$(pwd)/config.yml,dst=/app/config.yml,readonly \
    --mount type=bind,src=$(pwd)/db,dst=/app/db \
    -p 0.0.0.0:8080:80 \
    psyb0t/telegram-logger
```

### Run via docker compose

docker-compose.yml

```yaml
version: "3.8"
services:
  server:
    image: psyb0t/telegram-logger
    ports:
      - 8080:80
    volumes:
      - ./config.yml:/app/config.yml:ro
      - ./db:/app/db
    environment:
      - CONFIGFILE=/app/config.yml
      - LOGGER_LEVEL=debug
      - LOGGER_FORMAT=json
      - TELEGRAMBOT_TOKEN=TOKEN
      - TELEGRAMBOT_SUPERUSERCHATID=CHATID
      - STORAGE_TYPE=badgerDB
      - STORAGE_BADGERDB_DSN=/app/db
```

```
docker compose -f docker-compose.yml up
```

### Run via latest released binary

`download-latest.sh`

```bash
#!/usr/bin/env bash
owner=psyb0t
repo=telegram-logger
asset_name=telegram-logger-linux-amd64

echo "Looking up the latest release of $asset_name for github.com/$owner/$repo..."
releases=$(curl -s https://api.github.com/repos/$owner/$repo/releases)
latest_release=$(echo "$releases" | jq -r '.[0]')
asset_url=$(echo "$latest_release" | jq -r ".assets[] | select(.name == \"$asset_name\") | .browser_download_url")

echo "Downloading $asset_url..."
curl -s -L -o "$asset_name" "$asset_url"

chmod +x "$asset_name"
```

`run.sh`

```bash
export LISTENADDRESS=0.0.0.0:8080
export LOGGER_LEVEL=debug
export LOGGER_FORMAT=json
export TELEGRAMBOT_TOKEN=TOKEN
export TELEGRAMBOT_SUPERUSERCHATID=CHATID
export STORAGE_TYPE=badgerDB
export STORAGE_BADGERDB_DSN=./db

./telegram-logger-linux-amd64
```

## Interacting with the service

### curl

```
curl -X POST -H "Content-Type: application/json" -H "X-ID: abcdef" -d '{
  "caller": "myService",
  "time": "2022-03-11T12:34:56.789Z",
  "level": "warning",
  "message": "Something went wrong!",
  "error": "Error: XYZ",
  "requestID": "abc123",
  "traceID": "def456",
  "spanID": "ghi789",
  "data": {
    "someKey": "someValue",
    "anotherKey": 123
  }
}' http://localhost:8080/
```

## TODO:

- more config validation
- add healthcheck
- add trace id to logs
- create external wrapper package based on telegramBotMessageHandler
- build embeddable helper packages to interact with a deployed service
