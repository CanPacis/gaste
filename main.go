package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/teris-io/shortid"
)

type UrlInstance struct {
	FullUrl  string `json:"full_url"`
	ShortUrl string `json:"short_url"`
}

type UrlResponse struct {
	FullUrl    string `json:"full_url"`
	ShortUrl   string `json:"short_url"`
	Identifier string `json:"identifier"`
}

var default_domain = "https://gaste.news"

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/generate_url",
			Handler: func(c echo.Context) error {
				full_url := c.QueryParam("full_url")

				if len(full_url) > 0 {
					record, err := app.Dao().FindFirstRecordByData("urls", "full_url", full_url)

					if err != nil {
						id, err := shortid.Generate()

						if err != nil {
							return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate id"})
						}

						short_url := fmt.Sprintf("%s/%s", default_domain, id)
						collection, err := app.Dao().FindCollectionByNameOrId("urls")
						if err != nil {
							return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to find the collection"})
						}

						record := models.NewRecord(collection)
						record.Set("full_url", full_url)
						record.Set("short_url", short_url)
						record.Set("identifier", id)

						if err := app.Dao().SaveRecord(record); err != nil {
							return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save the record"})
						}

						return c.JSON(http.StatusOK, map[string]UrlResponse{"message": {ShortUrl: short_url, FullUrl: full_url, Identifier: id}})
					}

					return c.JSON(http.StatusOK, map[string]UrlResponse{"message": {ShortUrl: record.Get("short_url").(string), FullUrl: record.Get("full_url").(string), Identifier: record.Get("identifier").(string)}})
				}

				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Bad Request"})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
