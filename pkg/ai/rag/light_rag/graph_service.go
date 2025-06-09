package light_rag

type GraphService struct {
	cli *Client
}

func NewGraphService(cli *Client) *GraphService {
	return &GraphService{cli: cli}
}
