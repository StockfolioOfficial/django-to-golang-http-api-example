package article

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/supporter"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const (
	tag = "[ARTICLE] "
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

func decodeCursor(encodedTime string) (time.Time, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(data)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

func encodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

type view struct {
	supporter.APIView
}

func (self *view) Get(ctx echo.Context) error {
	var data struct {
		Num int `query:"num"`
		Cursor string `query:"cursor"`
	}

	err := ctx.Bind(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, supporter.ErrorResponse(err))
	}

	if data.Num == 0 {
		data.Num = 10
	}

	cursor, err := decodeCursor(data.Cursor)
	if err != nil && data.Cursor != "" {
		log.WithError(err).Error(tag, "cursor decode failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}

	var list []Model
	err = object.Preload("Author").
		Limit(data.Num).
		Where("created_at > ?", cursor).
		Find(&list).Error
	if err != nil {
		log.WithError(err).Error(tag, "object get list failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}

	var nextCursor string
	if l := len(list); l == data.Num {
		nextCursor = encodeCursor(time.Time(list[l-1].CreatedAt))
	}

	ctx.Response().Header().Set("X-Cursor", nextCursor)

	var response = make([]echo.Map, len(list))
	for i := range list {
		var article = list[i]
		response[i] = echo.Map{
			"id":      article.ID,
			"title":   article.Title,
			"content": article.Content,
			"author": echo.Map{
				"id":         article.Author.ID,
				"name":       article.Author.Name,
				"updated_at": article.Author.UpdatedAt,
				"created_at": article.Author.CreatedAt,
			},
			"updated_at": article.UpdatedAt,
			"created_at": article.CreatedAt,
		}
	}
	return ctx.JSON(http.StatusOK, response)
}

func (self *view) Post(ctx echo.Context) error {
	var data struct {
		Title string `json:"title" validate:"required,min=2,max=45"`
		Content string `json:"content" validate:"required,min=4"`
		Author struct {
			ID int64 `json:"id" validate:"required"`
		} `json:"author" validate:"required,dive"`
	}

	err := ctx.Bind(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, supporter.ErrorResponse(err))
	}

	err = ctx.Validate(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, supporter.ErrorResponse(err))
	}

	var exists Model
	err = object.Where("`title` = ?", data.Title).First(&exists).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithError(err).Error(tag, "object get first failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	} else if exists != (Model{}) {
		return ctx.JSON(http.StatusConflict, supporter.ErrorResponseMessage("already exists"))
	}

	var article = Model{
		Title:     data.Title,
		Content:   data.Content,
		AuthorID:  data.Author.ID,
		UpdatedAt: supporter.Time(time.Now()),
		CreatedAt: supporter.Time(time.Now()),
	}

	err = object.First(&article.Author, data.Author.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusBadRequest, supporter.ErrorResponseMessage("not exists author id"))
		}

		log.WithError(err).Error(tag, "author get one by id failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}

	err = object.Create(&article).Error
	if err != nil {
		log.WithError(err).Error(tag, "object create failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}


	return ctx.JSON(http.StatusCreated, echo.Map{
		"id":      article.ID,
		"title":   article.Title,
		"content": article.Content,
		"author": echo.Map{
			"id":         article.Author.ID,
			"name":       article.Author.Name,
			"updated_at": article.Author.UpdatedAt,
			"created_at": article.Author.CreatedAt,
		},
		"updated_at": article.UpdatedAt,
		"created_at": article.CreatedAt,
	})
}

func articleAsView() supporter.View {
	return &view{}
}

type detailView struct {
	supporter.APIView
}

func (self *detailView) Get(ctx echo.Context) error {
	var data struct {
		ID int64 `param:"article_id" validate:"required,ne=0"`
	}

	err := ctx.Bind(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, supporter.ErrorResponse(err))
	}

	err = ctx.Validate(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, supporter.ErrorResponse(err))
	}

	var article Model
	err = object.Preload("Author").
		First(&article, data.ID).Error
	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, echo.Map{
			"id":      article.ID,
			"title":   article.Title,
			"content": article.Content,
			"author": echo.Map{
				"id":         article.Author.ID,
				"name":       article.Author.Name,
				"updated_at": article.Author.UpdatedAt,
				"created_at": article.Author.CreatedAt,
			},
			"updated_at": article.UpdatedAt,
			"created_at": article.CreatedAt,
		})
	case gorm.ErrRecordNotFound:
		return ctx.JSON(http.StatusNotFound, supporter.ErrorResponseMessage("not found"))
	default:
		log.WithError(err).Error(tag, "object get one by id failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}
}

func (self *detailView) Delete(ctx echo.Context) error {
	var data struct {
		ID int64 `param:"article_id" validate:"required,ne=0"`
	}

	err := ctx.Bind(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, supporter.ErrorResponse(err))
	}

	err = ctx.Validate(&data)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, supporter.ErrorResponse(err))
	}
	var article Model
	err = object.First(&article, data.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, supporter.ErrorResponseMessage("not found"))
		}

		log.WithError(err).Error(tag, "object get one by id failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}

	err = object.Delete(&article).Error
	if err != nil {
		log.WithError(err).Error(tag, "object get one by id failed")
		return ctx.JSON(http.StatusInternalServerError, supporter.ErrorResponseMessage("server internal error"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func articleDetailAsView() supporter.View {
	return &detailView{}
}
