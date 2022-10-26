package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Exam/api_gateway/api/handlers/models"
	"github.com/Exam/api_gateway/genproto/post"
	pbp "github.com/Exam/api_gateway/genproto/post"
	l "github.com/Exam/api_gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreatePost
// @Summary creation of posts
// @Description creating posts
// @Tags Post
// @Accept json
// @Produce json
// @Param post body post.PostReq true "post"
// @Success 201 {object} post.PostResp
// @Failure 400 "ErrorResponse"
// @Router /v1/post [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		postBody    post.PostReq
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&postBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().CreatePost(ctx, &postBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create customer", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPostById
// @Summary gets the post info
// @Description getting post info
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} post.GetPostResp
// @Failure 400 "ErrorResponse"
// @Router /v1/post/{id} [get]
func (h *handlerV1) GetPost(c *gin.Context) {
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

	responseCus, err := h.serviceManager.PostService().GetPost(ctx, &pbp.GetPostReq{
		Id: Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, responseCus)
}

// PostUpdate
// @Summary updates posts
// @Description updating posts
// @Tags Post
// @Accept json
// @Produce json
// @Param postbody body post.Post true "Update Post"
// @Success 201 {object} post.Post
// @Failure 400 "ErrorResponse"
// @Router /v1/posts [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		postBody    models.Post
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseEnumNumbers = true

	err := c.ShouldBindJSON(&postBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.PostService().UpdatePost(ctx, &pbp.Post{
		Name:        postBody.Name,
		Description: postBody.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary gets all posts
// @Description getting all posts
// @Tags Post
// @Accept json
// @Produce json
// @Success 200 {object} post.ListPostsResp
// @Router /v1/post/allposts [get]
func (h *handlerV1) ListPosts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().ListPosts(ctx, &post.Empty{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to collect all posts", l.Error(err))
	}
	c.JSON(http.StatusOK, response)
}

// DeletePost
// @Summary deletes post
// @Description deleting of post
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int true "delete Post"
// @Success 200 "success"
// @Failure 400 "ErrorResponse"
// @Router /v1/posts/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {

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

	customerResp, err := h.serviceManager.PostService().DeletePost(ctx, &pbp.Id{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
	}

	c.JSON(http.StatusCreated, customerResp)
}
