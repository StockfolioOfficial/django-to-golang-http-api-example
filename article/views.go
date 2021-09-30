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

type responseDTO struct {
	ID        int64             `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Author    authorResponseDTO `json:"author"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
}

type authorResponseDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func toDTO(src Model) responseDTO {
	return responseDTO{
		ID:      src.ID,
		Title:   src.Title,
		Content: src.Content,
		Author: authorResponseDTO{
			ID:        src.Author.ID,
			Name:      src.Author.Name,
			UpdatedAt: time.Time(src.Author.UpdatedAt),
			CreatedAt: time.Time(src.Author.CreatedAt),
		},
		UpdatedAt: time.Time(src.UpdatedAt),
		CreatedAt: time.Time(src.CreatedAt),
	}
}

func toDTOList(list []Model) (res []responseDTO) {
	res = make([]responseDTO, len(list))
	for i := range res {
		res[i] = toDTO(list[i])
	}

	return
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
	return ctx.JSON(http.StatusOK, toDTOList(list))
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

	return ctx.JSON(http.StatusCreated, toDTO(article))
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
		return ctx.JSON(http.StatusOK, toDTO(article))
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


//func getStatusCode(err error) int {
//	if err == nil {
//		return http.StatusOK
//	}
//
//	log.Error(err)
//	switch err {
//	case domain.ErrInternalServerError:
//		return http.StatusInternalServerError
//	case domain.ErrNotFound:
//		return http.StatusNotFound
//	case domain.ErrConflict:
//		return http.StatusConflict
//	default:
//		return http.StatusInternalServerError
//	}
//}