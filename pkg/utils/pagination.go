package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Pagination struct {
	Page         int64
	Limit        int64
	SortBy       string
	HasPrev      bool
	HasNext      bool
	TotalRecords int64
	TotalPages   int64
	CurrentPage  int64
	NextLink     string
	PrevLink     string
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func DefaultPaginationQueryParams(c *gin.Context) (int64, int64, string) {
	page := c.Query("page")
	limit := c.Query("limit")
	sortBy := c.Query("sort_by")
	var p, l int = 0, 0
	if page == "" {
		p = 1
	} else {
		p, _ = strconv.Atoi(page)
	}
	if limit == "" {
		l = 10
	} else {
		l, _ = strconv.Atoi(limit)
	}
	if sortBy == "" {
		sortBy = "new"
	}

	if l > 100 {
		l = 100
	}
	return int64(p), int64(l), sortBy
}

func BuildPagination(page int64, limit int64, sortBy string, total int64) Pagination {
	pagination := Pagination{}
	pagination.Page = page
	pagination.Limit = limit
	pagination.SortBy = sortBy
	pagination.TotalPages = total / limit
	pagination.TotalRecords = total
	pagination.CurrentPage = page

	if page < pagination.TotalPages {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if page == 0 {
		pagination.HasPrev = false
	} else {
		pagination.HasPrev = true
	}
	pagination.NextLink = "?page=" + strconv.Itoa(int(page)+1) + "&limit=" + strconv.Itoa(int(limit))
	pagination.PrevLink = "?page=" + strconv.Itoa(int(page)-1) + "&limit=" + strconv.Itoa(int(limit))
	return pagination
}

func UserNickName(user interface{}) string {
	if user == nil {
		return ""
	}
	nickname := user.(map[string]interface{})["nickname"].(string)
	return nickname
}
