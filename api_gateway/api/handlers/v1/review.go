package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Exam/api_gateway/genproto/review"
	pbp "github.com/Exam/api_gateway/genproto/review"
	l "github.com/Exam/api_gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateReview
// @Summary creation of reviews
// @Description creating reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param review body review.Review true "review"
// @Success 201 {object} review.Review
// @Failure 400 "ErrorResponse"
// @Router /v1/review [post]
func (h *handlerV1) CreateReview(c *gin.Context) {
	var (
		reviewBody  review.Review
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&reviewBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.ReviewService().CreateReview(ctx, &reviewBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create review", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetReviewById
// @Summary gets the review info
// @Description getting review info
// @Tags Review
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} review.Review
// @Failure 400 "ErrorResponse"
// @Router /v1/review/{id} [get]
func (h *handlerV1) GetReview(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	guid := c.Param("id")
	Id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.ReviewService().GetReview(ctx, &pbp.GetReviewReq{Id: Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}

// UpdateReview
// @Summary updates reviews
// @Description updating reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param review body review.Review true "rewiew"
// @Success 200 {object} review.Review
// @Failure 400 "ErrorResponse"
// @Router /v1/reviews [put]
func (h *handlerV1) UpdateReview(c *gin.Context) {
	var (
		reviewBody  review.Review
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseEnumNumbers = true

	err := c.ShouldBindJSON(&reviewBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.ReviewService().CreateReview(ctx, &reviewBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update review", l.Error(err))
	}
	c.JSON(http.StatusCreated, response)
}

// DeleteReview
// @Summary deletes reviews
// @Description deleting of reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param id path int true "delete review"
// @Success 200 "success"
// @Failure 400 "ErrorResponse"
// @Router /v1/reviews/{id} [delete]
func (h *handlerV1) DeleteReview(c *gin.Context) {

	var (
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseEnumNumbers = true
	guid := c.Param("id")

	id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.ReviewService().DeleteReview(ctx, &pbp.Id{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete review", l.Error(err))
	}

	c.JSON(http.StatusCreated, response)
}

// GetPostReview
// @Summary gets the post reviews
// @Description getting post reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} review.Review
// @Failure 400 "ErrorResponse"
// @Router /v1/review/post/{id} [get]
func (h *handlerV1) PostReview(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	guid := c.Param("id")
	Id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.ReviewService().GetPostReview(ctx, &pbp.GetReviewPost{PostId: Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}
