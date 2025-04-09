package dto

import (
	"inventory-service/pkg/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReadParamID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil

}

func ReadInt(ctx *gin.Context, key string, defaultValue int, v *validator.Validator) int {
	s := ctx.Query(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func ReadString(ctx *gin.Context, key, defaultValue string) string {
	s := ctx.Query(key)
	if s == ""{
		return defaultValue
	}
	return s

}
