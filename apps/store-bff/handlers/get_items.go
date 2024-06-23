package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/demo-service-discovery/apps/store-bff/pkg/cfg"
	"github.com/thoughtgears/demo-service-discovery/apps/store-bff/pkg/run_requests"
)

func GetItems(config *cfg.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count := ctx.Query("count")
		if count == "" {
			count = "30"
		}

		client, err := run_requests.NewClient(config, "item-api")
		if err != nil {
			log.Error().Err(err).Msg("Error creating client")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		resp, err := client.Do(ctx, "GET", fmt.Sprintf("/items?count=%s", count))
		if err != nil {
			log.Error().Err(err).Msg("Error making request")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var response interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Error().Err(err).Msg("Error decoding response")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, response)
	}

}
