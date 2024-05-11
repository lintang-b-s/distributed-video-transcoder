package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lintang/video-processing-worker/biz/domain"
	"lintang/video-processing-worker/config"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DkronAPI struct {
	BaseURL      string
	MyServiceURL string
}

func CreateDkronAPI(cfg *config.Config) *DkronAPI {
	return &DkronAPI{
		BaseURL:      cfg.Dkron.DkronURL,
		MyServiceURL: cfg.MyServiceURL,
	}
}

type JobReq struct {
	Name           string            `json:"name"`
	DisplayName    string            `json:"displayname"`
	Schedule       string            `json:"schedule"`
	Timezone       string            `json:"timezone"`
	Owner          string            `json:"owner"`
	OwnerEmail     string            `json:"owner_email"`
	Disabled       bool              `json:"disabled"`
	Concurrency    string            `json:"concurrency"`
	Executor       string            `json:"executor"`
	ExecutorConfig map[string]string `json:"executor_config"`
}

func (d *DkronAPI) AddJobUploadPlaylistToMinio(ctx context.Context, filename string) error {
	randomString := uuid.New().String()

	cronURL := "http://%s/api/v1/transcoder/upload_playlist"
	jobName := filename + randomString

	at := time.Now().Add(time.Duration(500) * time.Millisecond)
	payload, err := json.Marshal(JobReq{
		Name:        jobName,
		DisplayName: jobName,
		Schedule:    fmt.Sprintf("@at " + at.Format(time.RFC3339)),
		Timezone:    "Asia/Jakarta",
		Owner:       "lintang birda saputra",
		OwnerEmail:  "lintangbirdasaputra23@gmail.com",
		Disabled:    false,
		Concurrency: "allow",
		Executor:    "shell",
		ExecutorConfig: map[string]string{
			// "shell": "true",
			"command": `curl --location ` + cronURL + ` \
			--header 'Content-Type: application/json' \
			--data '{
				"filename": "` + filename + `"
			}'`,
		},
	})

	if err != nil {
		zap.L().Error("Marshal JSON", zap.Error(err), zap.String("filename", filename))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	req, err := http.NewRequest("POST", d.BaseURL, bytes.NewBuffer(payload))

	if err != nil {
		zap.L().Error("NewRequest ", zap.Error(err), zap.String("filename", filename))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("client.Do(req) ", zap.Error(err), zap.String("filename", filename))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	defer resp.Body.Close()
	return nil 
}
