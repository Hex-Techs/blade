package view

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetID(c *gin.Context) (uint, error) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}
