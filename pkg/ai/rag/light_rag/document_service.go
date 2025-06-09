package light_rag

import (
	"context"
	"github.com/spf13/afero"
	"io/fs"
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

type Statuses struct {
	PENDING    []*Document `json:"PENDING"`
	PROCESSED  []*Document `json:"PROCESSED"`
	PROCESSING []*Document `json:"PROCESSING"`
	FAILED     []*Document `json:"FAILED"`
}

type StatusesResponse struct {
	Statuses Statuses `json:"statuses"`
}

type DeleteByFileNameRequest struct {
	FileName string `json:"fileName"`
}

type DeleteByFileNameResponse struct {
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

func (s *DocumentService) Upload(ctx context.Context, fileName string, fs afero.Fs) error {
	return s.cli.uploadFile(ctx, "/documents/upload", fileName, fs)
}

func (s *DocumentService) DeleteByFileName(ctx context.Context, fileName string) (*DeleteByFileNameResponse, error) {
	request := &DeleteByFileNameRequest{
		FileName: fileName,
	}
	resp, err := s.cli.Post(ctx, "/documents/delete_by_filename", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res *DeleteByFileNameResponse
	if err = s.cli.getJsonData(ctx, resp, &res); err != nil {
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

func (s *DocumentService) GetStatuses(ctx context.Context) *StatusesResponse {
	return &StatusesResponse{Statuses: Statuses{}}
}
