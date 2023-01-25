package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	d "url-shorter/app/internal/controllers/v1/dto"
	"url-shorter/app/internal/domain/service/reduction"
)

var (
	api            = "/api/v1/"
	createShortUrl = "/createShortUrl/"
	getLongUrl     = "/getLongUrl/"
	redirect       = "/:hashed_link"
)

type ReductionHandler struct {
	service *reduction.ReductionService
}

func NewReductionHandler(service *reduction.ReductionService) *ReductionHandler {
	return &ReductionHandler{service: service}
}

func (r ReductionHandler) Register(router *gin.Engine) {
	v1 := router.Group(api)
	{
		v1.POST(createShortUrl, r.CreateShortUrl)
		v1.POST(getLongUrl, r.GetLongUrl)
	}

	router.GET(redirect, r.redirect)
}

// CreateShortUrl Create short url
// @Summary      Create short url
// @Description  Create short url
// @Produce      json
// @Accept      json
// @Param 		message body dto.CreateShortUrlDTO true "Create Short Url"
// @Success      201
// @Failure      400
// @Router      /api/v1/createShortUrl/ [post]
func (r ReductionHandler) CreateShortUrl(c *gin.Context) {
	var createShortUrlDTO d.CreateShortUrlDTO

	if err := c.BindJSON(&createShortUrlDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "failed",
			"msg":         err.Error(),
			"hashed_link": "",
			"short_url":   "",
		})

		return
	}

	url, err := r.service.CreateShortUrl(c.Request.Context(), reduction.CreateShortUrlDTO(createShortUrlDTO))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "failed",
			"msg":         err.Error(),
			"hashed_link": "",
			"short_url":   "",
		})

		return
	}

	hashedLink := strings.Split(url, "/")

	c.JSON(http.StatusCreated, gin.H{
		"status":      "success",
		"msg":         "",
		"hashed_link": hashedLink[len(hashedLink)-1],
		"short_url":   url,
	})
}

// GetLongUrl Get long url
// @Summary      Get long url
// @Description  Get long url
// @Produce      json
// @Accept      json
// @Param 		message body dto.GetLongUrlDTO true "Get Long Url"
// @Success      200
// @Failure      400
// @Router       /api/v1/getLongUrl/ [post]
func (r ReductionHandler) GetLongUrl(c *gin.Context) {
	var getLongUrlDTO d.GetLongUrlDTO

	if err := c.BindJSON(&getLongUrlDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":   "failed",
			"msg":      err.Error(),
			"long_url": "",
		})

		return
	}

	url, err := r.service.GetLongUrl(c.Request.Context(), getLongUrlDTO.HashedLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":   "failed",
			"msg":      err.Error(),
			"long_url": "",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"msg":      "",
		"long_url": url,
	})
}

// Redirect to original url
// @Summary      Redirect to original url
// @Description  Redirect to original url
// @Param		hashed_link	path string true "Hashed Link"
// @Success      302
// @Failure      400
// @Router       /{hashed_link} [get]
func (r ReductionHandler) redirect(c *gin.Context) {
	hashedLink := c.Param("hashed_link")

	url, err := r.service.GetLongUrl(c.Request.Context(), hashedLink)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte("Failed to redirect"))

		return
	}

	c.Redirect(http.StatusFound, url)
}
