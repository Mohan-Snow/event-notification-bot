# Scheduling bot

---
The primary objective of this project is to integrate the Telegram Bot API library and explore the interaction with Outlook calendar events using the Microsoft Outlook API. 

Tutorial how to get Telegram API Token: https://core.telegram.org/bots/tutorial

---
## Deployment:

1. Install [Task](https://taskfile.dev/installation/)
2. Create .env file and fill like [.env.example](./.env.example)
3. Run command `task deploy`
---

Branch and Commit messages naming rules:
Branch name:
```
<issue-number> - <Short issue name with dashes>
```
Commit message:
```
[<issue-number>] <Your commit message>
```

## Technologies Used

### Libraries and Dependencies

1. **Telegram Bot API**
    - **Library**: `github.com/go-telegram-bot-api/telegram-bot-api/v5`
    - **Version**: v5.5.1
    - **Description**: This library provides a wrapper for the Telegram Bot API, enabling us to interact with Telegram's messaging platform seamlessly and efficiently. It supports various bot functionalities such as sending messages, receiving updates, and handling different types of Telegram interactions.

2. **Environment Configuration**
    - **Library**: `github.com/kelseyhightower/envconfig`
    - **Version**: v1.4.0
    - **Description**: This library simplifies the process of managing environment variables in our application. It allows us to define, load, and use environment variables effortlessly, ensuring our application is configurable and adheres to the 12-factor app principles.

3. **PostgreSQL Driver**
    - **Library**: `github.com/lib/pq`
    - **Version**: v1.10.9
    - **Description**: The `pq` library is a pure Go driver for PostgreSQL, enabling our application to interact with PostgreSQL databases. It supports a wide range of PostgreSQL features, ensuring reliable and performant database operations.

4. **Structured Logging**
    - **Library**: `go.uber.org/zap`
    - **Version**: v1.27.0
    - **Description**: Zap is a high-performance, structured logging library for Go. It allows us to log detailed information in a structured format, making it easier to query and analyze logs. Zap provides a balance between performance and flexibility, supporting both human-readable and machine-readable log formats.
