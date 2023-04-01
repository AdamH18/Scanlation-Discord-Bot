# Deploying Scanlation-Discord-Bot using docker

## [ ! ] This documentation is incomplete due to the nature of the project. Use with caution.

### 1. Building Image

1. Install [docker environment](https://docs.docker.com/engine/install/) with docker compose
2. Clone this repository

```
$ git clone https://github.com/AdamH18/Scanlation-Discord-Bot.git && cd Scanlation-Discord-Bot
```

3. Build image

```
$ docker build . -t scanlation-discord-bot:latest
```

4. Run image using docker compose

```
$ cd docker
$ docker compose up -d
```

### Configuration

This app require specific value on configuration file in config.json

```
{
    "Token" : "token",
    "RemoveCommands" : "False",
    "DatabaseFile" : "sqlite.db",
    "DatabaseBackupChannel" : "ChannelID"
}
```

| key | value |
|---|---|
| Token | [Your discord bot token](https://discord.com/developers/docs/reference#authentication)|
| RemoveCommands | True / False |
| DatabaseFile | database name |
| DatabaseBackupChannel | [ChannelID](https://discord.com/developers/docs/resources/channel) |

Make sure the absolute volume path is adjusted to your local system for Docker-specific deployment.

```
    volumes:
      - /path/config.json:/app/config.json 
      - /path/sqlite.db:/app/sqlite.db
```
