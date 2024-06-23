package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/demo-service-discovery/apps/items-api/db"
)

func GetItems(c *gin.Context) {
	count := 30
	queryCount := c.Query("count")
	if queryCount != "" {
		if parsedCount, err := strconv.Atoi(queryCount); err == nil {
			count = parsedCount
		}
	}

	items := db.GetItems(count)
	c.JSON(200, items)
}
