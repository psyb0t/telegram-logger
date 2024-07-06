# telegram-logger: Hack Your Logs

Welcome to `telegram-logger`, where we turn your logs into Telegram messages. Because logs should be as mobile as you are.

## Table of Contents

- [What's This?](#whats-this)
- [Configuration](#configuration)
- [HTTP API](#http-api)
- [Telegram Bot](#telegram-bot)
- [Running the Service](#running-the-service)
- [Interacting with the Service](#interacting-with-the-service)
- [Available Clients](#available-clients)
- [TODO](#todo)

## What's This?

`telegram-logger` is a server that receives log entries via HTTP POST requests and sends them to a Telegram chat instance. It's like having a secret agent for your logs.

## Configuration

Your config file is the key to the kingdom. Here's what you can tweak:

```yaml
listenAddress: 0.0.0.0:8080
logger:
  level: debug
  format: json
telegramBot:
  token: YOUR_SECRET_TOKEN_HERE
  superuserChatID: 38081130
storage:
  type: badgerDB
  badgerDB:
    dsn: /path/to/db/dir
```

Prefer environment variables? We've got you covered:

```bash
export LISTENADDRESS=0.0.0.0:8080
export LOGGER_LEVEL=debug
export LOGGER_FORMAT=json
export TELEGRAMBOT_TOKEN=YOUR_SECRET_TOKEN_HERE
export TELEGRAMBOT_SUPERUSERCHATID=38081130
export STORAGE_TYPE=badgerDB
export STORAGE_BADGERDB_DSN=/path/to/db/dir
```

## HTTP API

Send a POST request with your log entry. Don't forget your secret handshake (X-ID header).

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

## Telegram Bot

Our bot's got a few commands that you can throw at it:

- `/start`: Get your unique ID
- `/stop`: Go dark
- `/getAllUsers`: For the admins
- `/addUser`: Recruit new agents (admin only) - a user can also be a channel

Pro Tip: Adding a channel? Here's how:

1. Send a message to your target channel
2. Grab the channel ID from the message URL (e.g., https://t.me/c/2340157712/5)
3. Add the prefix: `-100`
4. Use the command: `/addUser -1002340157712`
5. Grab the ID that the bot sends to the channel(and maybe delete that message) and use it as your `X-ID` header when doing your HTTP request.

## Running the Service

### Docker

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

### Docker Compose

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

Fire it up:

```bash
docker compose -f docker-compose.yml up
```

### Binary

First, get the latest binary:

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

Then run it:

```bash
export LISTENADDRESS=0.0.0.0:8080
export LOGGER_LEVEL=debug
export LOGGER_FORMAT=json
export TELEGRAMBOT_TOKEN=YOUR_SECRET_TOKEN_HERE
export TELEGRAMBOT_SUPERUSERCHATID=YOUR_CHAT_ID_HERE
export STORAGE_TYPE=badgerDB
export STORAGE_BADGERDB_DSN=./db

./telegram-logger-linux-amd64
```

## Interacting with the Service

Time to put it to the test:

```bash
curl -X POST -H "Content-Type: application/json" -H "X-ID: YOUR_SECRET_ID" -d '{
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

## Available Clients

- https://github.com/psyb0t/py-telegram-logger-client: For the Pythonistas (module and executable)

More clients in the pipeline. Stay tuned.

## TODO

- More config validation
- Add healthcheck
- Add trace id to logs
- Create external wrapper package based on telegramBotMessageHandler
- Build embeddable helper packages to interact with a deployed service
- Fix linting
- SQLITE support

Now go forth and hack those logs! 🖥️
