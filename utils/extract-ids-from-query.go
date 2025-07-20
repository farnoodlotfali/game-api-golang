package utils

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

type HasID interface {
	GetID() int64
}

func ExtractIDsFromQuery[T HasID](ctx *gin.Context, queryParam string) ([]int64, error) {
	param := ctx.DefaultQuery(queryParam, "")
	if param == "" {
		return nil, nil
	}

	decodedParam, err := url.QueryUnescape(param)
	if err != nil {
		return nil, fmt.Errorf("invalid encoding for %s", queryParam)
	}

	var items []T
	if err := json.Unmarshal([]byte(decodedParam), &items); err != nil {
		return nil, fmt.Errorf("invalid format for %s", queryParam)
	}

	var ids []int64
	for _, item := range items {
		ids = append(ids, item.GetID())
	}
	return ids, nil
}
