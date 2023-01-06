# telegram-logger

[![codecov](https://codecov.io/gh/psyb0t/telegram-logger/branch/master/graph/badge.svg?token=QG0NA3QE7I)](https://codecov.io/gh/psyb0t/telegram-logger)
[![goreportcard](https://goreportcard.com/badge/github.com/psyb0t/telegram-logger)](https://goreportcard.com/report/github.com/psyb0t/telegram-logger)
[![pipeline](https://github.com/psyb0t/telegram-logger/actions/workflows/pipeline.yml/badge.svg)](https://github.com/psyb0t/telegram-logger/actions/workflows/pipeline.yml)

This service is a server that receives log entries via HTTP POST requests and sends them to a Telegram bot. The log entries are expected to be in the form of a JSON object with various fields representing log data and metadata.

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

The configuration file can be specified using the `CONFIGFILE` environment variable or, if not set, the default file `./config.yml` will be used. Environment variables can also be used to override values in the configuration file (e.g. `LOGGER_LEVEL` or `TELEGRAMBOT_TOKEN`).

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

## HTTP API

The service listens for HTTP POST requests to the root path and sends the request body as a message via the Telegram bot to the user associated with the request (determined by the value of the `X-ID` header). The response body is a JSON object with the following fields:

- error: an error message if an error occurred, otherwise empty.
- message: a message indicating the result of the request.

```json
{
  "error": "Something went wrong!",
  "message": "Your request could not be processed"
}
```

## Telegram Bot

The service has a number of commands that can be sent via the Telegram bot by a user. These commands are handled by the app and invoke the corresponding command handler function. The following commands are supported:

- /start: generates a unique ID for the user and sends it to the user in a welcome message. The user can then use this ID in the X-ID header of HTTP POST requests to the app to authenticate themselves.
- /stop: stops sending log entries to the user and removes the user's ID from the database.
- /getAllUsers (superadmin only): retrieves a list of all users from the database and sends it to the user via the Telegram bot.

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
