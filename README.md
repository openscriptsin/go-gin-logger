# File Based Logger for Gin in Golang
This is a simple guide to set up a file-based logger for a Gin server in Golang. This logger provides four methods to write logs to respective files: info.log, debug.log, warn.log, and error.log.

## Requirements
- Go installed on your system
- Gin framework [github.com/gin-gonic/gin] installed in your project

## Installation
```
go get github.com/openscriptsin/go-logger
```

## Usage

```
import (
    logger github.com/openscriptsin/go-logger
)

var logConfig logger.LoggerConfig = logger.LoggerConfig{
	Env:    "LOCAL",  // STAGE or PROD
	LogDir: "logs",
}

func main() {

    // creating instance of Logger
	appLogger := logger.NewGinLogger(
        logConfig, logger.ContextField("X-Request-Id"), logger.ContextField("X-User-Id")
    )()

    // ctx *gin.Context
    appLogger.Info(ctx, "calling middleware status called")
}

```

## Dependencies
- [Logrus]
- [Gin]


## Example 
- [go-gin-poc]

## License
MIT

## Acknowledgments
- Hat tip to the contributors of Gin and Logrus.






[//]: # (Links)
[github.com/gin-gonic/gin]: <https:github.com/gin-gonic/gin>
[go-gin-poc]: <https://github.com/amiransari27/go-gin-poc>
[Logrus]: <https://pkg.go.dev/github.com/sirupsen/logrus>
[Gin]: <https://pkg.go.dev/github.com/gin-gonic/gin>