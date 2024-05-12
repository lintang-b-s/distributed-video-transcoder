package router

import (
	"context"
	"errors"
	"lintang/video-processing-worker/biz/domain"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type TranscoderService interface {
	Transcode(ctx context.Context, filename string, resolution Resolution) error
	GenerateDASHPlaylist(ctx context.Context, filename string) error
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
		tH := root.Group("/transcoder")
		{
			tH.POST("/upload_playlist", handler.UploadPlaylistToMinio)
			tH.POST("/transcode", handler.Transcode)
		}
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}
type cronResp struct {
	Message string `json:"string"`
}

type Resolution string

const (
	Res240p  Resolution = "240p"
	Res360p  Resolution = "360p"
	Res480p  Resolution = "480p"
	Res720p  Resolution = "720p"
	Res1080p Resolution = "1080p"
)

type transcodeReq struct {
	Filename   string     `json:"filename" vd:"len($) <200; msg:'panjang filename harus kurang dari 200'"`
	Resolution Resolution `json:"resolution" `
}

func (h *TranscoderHandler) Transcode(ctx context.Context, c *app.RequestContext) {
	var req transcodeReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	err = h.ts.Transcode(ctx, req.Filename, req.Resolution)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

type uploadPlaylistReq struct {
	Filename string `query:"filename" vd:"len($) <200; msg:'panjang filename harus kurang dari 200'"`
}

func (h *TranscoderHandler) UploadPlaylistToMinio(ctx context.Context, c *app.RequestContext) {
	var req uploadPlaylistReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	err = h.ts.GenerateDASHPlaylist(ctx, req.Filename)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, cronResp{Message: "ok"})
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
