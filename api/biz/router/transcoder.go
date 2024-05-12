package router

import (
	"context"
	"errors"
	"lintang/video-transcoder-api/biz/domain"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type TranscoderService interface {
	CreatePresignedURLForUpload(ctx context.Context, filename string) (string, error)
	GetAllVideosMetadata(ctx context.Context) ([]domain.VideoMetadata, error)
}

type TranscoderHandler struct {
	ts TranscoderService
}

func TranscoderRouter(r *server.Hertz, ts TranscoderService) {
	handler := &TranscoderHandler{
		ts,
	}

	root := r.Group("/api/v1")
	{
		tH := root.Group("/tenflix")
		{
			tH.POST("/upload", handler.CreatePresignedURLForUpload)
			tH.GET("/", handler.GetAll)
		}
	}
}

type ResponseError struct {
	Message string `json:"message"`
}

type createPresignedURLReq struct {
	Filename string `json:"filename,required"`
}

type createPresignedURLResp struct {
	PresignedURL string `json:"presigned_url"`
}

type getAllMetadatasResp struct {
	Metadatas []domain.VideoMetadata `json:"metadatas"`
}

func (h *TranscoderHandler) GetAll(ctx context.Context, c *app.RequestContext) {
	metadatas, err := h.ts.GetAllVideosMetadata(ctx)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getAllMetadatasResp{metadatas})
}

func (h *TranscoderHandler) CreatePresignedURLForUpload(ctx context.Context, c *app.RequestContext) {
	var req createPresignedURLReq
	err := c.BindAndValidate(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	presignedURL, err := h.ts.CreatePresignedURLForUpload(ctx, req.Filename)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, createPresignedURLResp{PresignedURL: presignedURL})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ierr *domain.Error
	if !errors.As(err, &ierr) {
		return http.StatusInternalServerError
	} else {
		switch ierr.Code() {
		case domain.ErrInternalServerError:
			return http.StatusInternalServerError
		case domain.ErrNotFound:
			return http.StatusNotFound
		case domain.ErrConflict:
			return http.StatusConflict
		case domain.ErrBadParamInput:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

}
