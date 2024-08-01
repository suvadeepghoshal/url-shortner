package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CommonResponse struct {
	Time time.Time
}

func main() {
	log.Infof("inside main()")
	commonResponse := CommonResponse{
		Time: time.Now(),
	}
	e := echo.New()
	// have new init handler
	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, commonResponse)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
