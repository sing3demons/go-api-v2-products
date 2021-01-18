package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

func pagingResource(ctx *gin.Context, query *gorm.DB, records interface{}) *pagingResult {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "12"))

	ch := make(chan int)
	go countRecords(query, records, ch)

	offset := (page - 1) * limit
	query.Offset(offset).Limit(limit).Find(records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))
	// 5. Find nextPage
	var nextPage int
	if nextPage == totalPage {
		nextPage = totalPage
	} else {
		nextPage = totalPage + 1
	}
	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		Count:     count,
		TotalPage: totalPage,
	}

}

func countRecords(query *gorm.DB, records interface{}, ch chan int) {
	var count int64
	query.Model(records).Count(&count)

	ch <- int(count)
}
