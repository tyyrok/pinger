## Pinger Bot with Telegram Notifications

A lightweight Go-based bot that monitors the availability of specified URLs and sends Telegram notifications when a service goes down. The Docker image is compact—only ~27 MB.

### Repository Includes:

- Go service code
- docker-compose.yml for running the bot in a container
- Shell script for easy startup

### Setup Instructions
1. Create a settings.json file in the root directory with the following structure:
```
{
  "hosts": [
    {
      "My Service": "http://example.com"
    }
  ],
  "users": {"123456789": false}
}
```
- Replace "My Service" and "http://example.com" with your actual service name and URL.

2. Create a .env file with your bot's token:
`BOT_TOKEN=your_telegram_bot_token`

### Running the Bot

1. Start services:`./start.sh` or `docker-compose up --build -d`
2. Subscribe to bot notifications by sending `/start` command to your bot
3. Unsubscribe from bot sending `/stop` commnand to your bot
4. Get instant status of added hosts by sending `/status`
