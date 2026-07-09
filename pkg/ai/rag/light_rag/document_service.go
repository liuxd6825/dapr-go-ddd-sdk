package light_rag

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"io/fs"
	"strings"
)

type Document struct {
	ContentLength  int    `json:"content_length"`
	ContentSummary string `json:"content_summary"`
	CreatedAt      string `json:"created_at"`
	FilePath       string `json:"file_path"`
	Id             string `json:"id"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at"`
}

type TrackDocument struct {
	Documents []*Document `json:"documents"`
	TrackId   string      `json:"track_id"`
}

type Statuses struct {
	PENDING    []*Document `json:"PENDING"`
	PROCESSED  []*Document `json:"PROCESSED"`
	PROCESSING []*Document `json:"PROCESSING"`
	FAILED     []*Document `json:"FAILED"`
	COMPLETED  []*Document `json:"COMPLETED"`
	ANALYZING  []*Document `json:"ANALYZING"`
}

type StatusesResponse struct {
	Statuses Statuses `json:"statuses"`
}

type DeleteByIdsRequest struct {
	DocIds         []string `json:"doc_ids"`
	DeleteFile     bool     `json:"delete_file"`
	DeleteLlmCache bool     `json:"delete_llm_cache"`
}

type DeleteByIdsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type DocumentService struct {
	cli *Client
}

func NewDocumentService(cli *Client) *DocumentService {
	return &DocumentService{cli: cli}
}

func (s *DocumentService) Scan() {

}

func (s *DocumentService) Upload(ctx context.Context, fileName string, fs afero.Fs) (string, error) {
	return s.cli.uploadFile(ctx, "/documents/upload", fileName, fs)
}

func (s *DocumentService) Delete(ctx context.Context, ids []string) (*DeleteByIdsResponse, error) {
	request := &DeleteByIdsRequest{
		DocIds:         ids,
		DeleteFile:     true,
		DeleteLlmCache: true,
	}
	resp, err := s.cli.Delete(ctx, "/documents/delete_document", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := &DeleteByIdsResponse{}
	if err = s.cli.getJsonData(ctx, resp, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *DocumentService) InsertText(ctx context.Context, text string) {

}

func (s *DocumentService) InsertTexts(ctx context.Context) {

}

func (s *DocumentService) InsertFile(ctx context.Context, fs fs.File) {

}

func (s *DocumentService) InsertFileBatch(ctx context.Context, fsList []fs.File) {

}

func (s *DocumentService) Clear(ctx context.Context) {

}

func (s *DocumentService) GetStatuses(ctx context.Context, trackId string) (*StatusesResponse, error) {
	resp, err := s.cli.Get(ctx, fmt.Sprintf("/documents/track_status/%s", trackId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := &TrackDocument{}
	if err = s.cli.getJsonData(ctx, resp, res); err != nil {
		return nil, err
	}
	ss := Statuses{
		PENDING:    []*Document{},
		PROCESSED:  []*Document{},
		PROCESSING: []*Document{},
		FAILED:     []*Document{},
		COMPLETED:  []*Document{},
		ANALYZING:  []*Document{},
	}
	for _, o := range res.Documents {
		if strings.ToUpper(o.Status) == "PENDING" {
			ss.PENDING = append(ss.PENDING, o)
		} else if strings.ToUpper(o.Status) == "PROCESSED" {
			ss.PROCESSED = append(ss.PROCESSED, o)
		} else if strings.ToUpper(o.Status) == "PROCESSING" {
			ss.PROCESSING = append(ss.PROCESSING, o)
		} else if strings.ToUpper(o.Status) == "FAILED" {
			ss.FAILED = append(ss.FAILED, o)
		} else if strings.ToUpper(o.Status) == "COMPLETED" {
			ss.COMPLETED = append(ss.COMPLETED, o)
		} else if strings.ToUpper(o.Status) == "ANALYZING" {
			ss.ANALYZING = append(ss.ANALYZING, o)
		} else {
			print("不知道的状态", o.Status)
		}
	}
	return &StatusesResponse{Statuses: ss}, nil
}
