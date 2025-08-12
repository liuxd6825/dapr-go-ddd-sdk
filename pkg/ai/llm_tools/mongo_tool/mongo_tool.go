package mongo_tool

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type Region string
type SafeSearch string
type TimeRange string

const (
	// Regions settings
	RegionUS Region = "en-US"
	RegionGB Region = "en-GB"
	RegionCA Region = "en-CA"
	RegionAU Region = "en-AU"
	RegionDE Region = "de-DE"
	RegionFR Region = "fr-FR"
	RegionCN Region = "zh-CN"
	RegionHK Region = "zh-HK"
	RegionTW Region = "zh-TW"
	RegionJP Region = "ja-JP"
	RegionKR Region = "ko-KR"

	// SafeSearch settings
	SafeSearchOff      SafeSearch = "Off"
	SafeSearchModerate SafeSearch = "Moderate"
	SafeSearchStrict   SafeSearch = "Strict"

	// TimeRange settings
	TimeRangeDay   TimeRange = "Day"
	TimeRangeWeek  TimeRange = "Week"
	TimeRangeMonth TimeRange = "Month"
)

// Config represents the Bing search tool configuration.
type Config struct {
	// Eino tool settings
	ToolName string `json:"tool_name"` // optional, default is "bing_search"
	ToolDesc string `json:"tool_desc"` // optional, default is "search web for information by bing"

	AppName        string `json:"app_name"`
	Addr           string `json:"addr"`
	RsName         string `json:"rs_name"`
	UserName       string `json:"user_name"`
	Password       string `json:"password"`
	AutoDB         string `json:"auto_db"`
	AuthMechanism  string `json:"auth_mechanism"`
	ConnectTimeout int    `json:"connect_timeout"`
	CollName       string `json:"coll_name"`

	// Bing search settings
	// APIKey The API key is required to access the Bing Web Search API.
	APIKey string `json:"api_key"`

	// Region specifies the Bing search region and is used to customize the search results for a specific country or language.
	// Optional, default: ""
	Region Region `json:"region"`

	// MaxResults specifies the maximum number of search results to return.
	// Optional, default: 10
	MaxResults int `json:"max_results"`

	// SafeSearch specifies the Bing search safe search setting.
	// Optional, default: SafeSearchModerate
	SafeSearch SafeSearch `json:"safe_search"`

	// TimeRange specifies the Bing search time range.
	// Optional, default: ""
	TimeRange TimeRange `json:"time_range"`

	// Bing client settings
	// Headers specifies custom HTTP headers to be sent with each request.
	// Common headers like "User-Agent" can be set here.
	// Optional, default: map[string]string{}
	// Example:
	//   Headers: map[string]string{
	//     "User-Agent": "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
	//     "Accept-Language": "en-US",
	//   }
	Headers map[string]string `json:"headers"`

	// Timeout specifies the maximum duration for a single request.
	// Optional, default: 30 * time.Second
	// Example: 5 * time.Second
	Timeout time.Duration `json:"timeout"`

	// ProxyURL specifies the proxy server URL for all requests.
	// Supports HTTP, HTTPS, and SOCKS5 proxies.
	// Optional, default: ""
	// Example values:
	//   - "http://proxy.example.com:8080"
	//   - "socks5://localhost:1080"
	//   - "tb" (special alias for Tor Browser)
	ProxyURL string `json:"proxy_url"`

	// Cache enables in-memory caching of search results.
	// When enabled, identical search requests will return cached results
	// for improved performance. Cache entries expire after 5 minutes.
	// Optional, default: 0 (disabled)
	// Example: 5 * time.Minute
	Cache time.Duration `json:"cache"`

	// MaxRetries specifies the maximum number of retry attempts for failed requests.
	// Optional, default: 3
	MaxRetries int `json:"max_retries"`
}

// NewTool creates a new Bing search tool instance.
func NewTool(ctx context.Context, config *Config) (tool.InvokableTool, error) {
	bing, err := newMongoQuery(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create bing search tool: %w", err)
	}

	searchTool, err := utils.InferTool(config.ToolName, config.ToolDesc, bing.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to infer tool: %w", err)
	}

	return searchTool, nil
}

// validate validates the Bing search tool configuration.
func (c *Config) validate() error {
	// Set default values
	if c.ToolName == "" {
		c.ToolName = "mongo_query"
	}

	if c.ToolDesc == "" {
		c.ToolDesc = "query MongoDB for information"
	}

	// Validate required fields
	if c.APIKey == "" {
		//return errors.New("bing search tool config is missing API key")
	}

	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}

	c.Headers["Ocp-Apim-Subscription-Key"] = c.APIKey

	return nil
}

// mongoQuery represents the Bing search tool.
type mongoQuery struct {
	config *Config
	client *mongo.Client
	coll   *mongo.Collection
}

// newMongoQuery creates a new Bing search client.
func newMongoQuery(config *Config) (*mongoQuery, error) {
	if config == nil {
		return nil, errors.New("bing search tool config is required")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	connectTimeout := time.Duration(config.ConnectTimeout)
	opt := &options.ClientOptions{
		AppName:        &config.AppName,
		ConnectTimeout: &connectTimeout,
		Auth: &options.Credential{
			Username:      config.UserName,
			Password:      config.Password,
			AuthSource:    config.AutoDB,
			AuthMechanism: config.AuthMechanism,
		},
	}

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &mongoQuery{
		config: config,
		client: client,
	}, nil
}

type QueryRequest struct {
	Query  string `json:"query" jsonschema:"description=The query to search the web for"`
	Offset int    `json:"page" jsonschema:"description=The index of the first result to return, default is 0"`
}

type QueryResponse struct {
	Results []map[string]any `json:"results" jsonschema:"description=The results of the search"`
}

// Query searches the web for information.
func (s *mongoQuery) Query(ctx context.Context, request *QueryRequest) (response *QueryResponse, err error) {
	// Search the web for information
	filter := bson.M{}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []map[string]any
	err = cur.All(ctx, &list)
	if err != nil {
		return nil, err
	}

	return &QueryResponse{
		Results: list,
	}, nil
}
