package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReadIDParam(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, nil
	}
	return id, nil
}
