package gateway

import (
	"bnsportal/gateway/respond"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (gw *Gateway) GetAddressNames(ctx *gin.Context) {
	address := ctx.Param("address")

	items, err := gw.Storage.GetAddressNames(strings.ToLower(address))

	if err != nil {
		wrapHTTPRespond(ctx, nil, err)
		return
	}
	result := []respond.NameInfo{}
	for _, item := range items {
		result = append(result, respond.NameInfo{
			Name:              item.Name,
			Owner:             item.Owner,
			ID:                item.ID,
			RegisteredAtBlock: item.RegisteredAtBlock,
		})
	}
	wrapHTTPRespond(ctx, result, err)
}

func (gw *Gateway) GetNames(ctx *gin.Context) {
	offset, _ := strconv.ParseInt(ctx.Query("offset"), 10, 64)

	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

	if limit == 0 {
		limit = 100
	}
	if limit > 100 {
		wrapHTTPRespond(ctx, nil, errors.New("limit should be less than 100"))
		return
	}
	filter := ctx.Query("filter")
	filterMap := make(map[string]interface{})
	if filter != "" {
		err := json.Unmarshal([]byte(filter), &filterMap)
		if err != nil {
			wrapHTTPRespond(ctx, nil, err)
			return
		}
	}
	items, err := gw.Storage.FilterName(limit, offset, filterMap)
	if err != nil {
		wrapHTTPRespond(ctx, nil, err)
		return
	}
	result := []respond.NameInfo{}
	for _, item := range items {
		result = append(result, respond.NameInfo{
			Name:              item.Name,
			Owner:             item.Owner,
			ID:                item.ID,
			RegisteredAtBlock: item.RegisteredAtBlock,
		})
	}
	wrapHTTPRespond(ctx, result, err)
}

func (gw *Gateway) CheckNameAvailable(ctx *gin.Context) {
	name := ctx.Param("name")

	available, err := gw.Storage.CheckNameAvailable(name)

	if err != nil {
		wrapHTTPRespond(ctx, nil, err)
		return
	}
	wrapHTTPRespond(ctx, available, err)
}

func (gw *Gateway) GetName(ctx *gin.Context) {
	name := ctx.Param("name")

	nameInfo, err := gw.Storage.GetNameInfo(name)

	if err != nil {
		wrapHTTPRespond(ctx, nil, err)
		return
	}
	result := respond.NameInfo{}
	result.Name = nameInfo.Name
	result.Owner = nameInfo.Owner
	result.ID = nameInfo.ID
	result.RegisteredAtBlock = nameInfo.RegisteredAtBlock
	wrapHTTPRespond(ctx, result, err)
}
