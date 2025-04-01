package swagger3

import (
	"strings"
)

type Swagger struct {
	Openapi    string      `json:"openapi"`
	Servers    []*Server   `json:"servers,omitempty"`
	Info       *Info       `json:"info,omitempty"`
	Tags       []*Tag      `json:"tags,omitempty"`
	Paths      Paths       `json:"paths,omitempty"`
	Components *Components `json:"components,omitempty"`
}

type HttpMethod string

const (
	HttpMethodGet    HttpMethod = "get"
	HttpMethodPost   HttpMethod = "post"
	HttpMethodPut    HttpMethod = "put"
	HttpMethodPatch  HttpMethod = "patch"
	HttpMethodDelete HttpMethod = "delete"
)

func NewSwagger() *Swagger {
	return &Swagger{
		Openapi: "3.0.1",
		Servers: []*Server{},
		Info:    &Info{},
		Tags:    []*Tag{},
		Paths:   Paths{},
		Components: &Components{
			Schemas: GetComponents(),
		},
	}
}

func (s *Swagger) Filter(resource string) *Swagger {
	if resource == "" {
		return s
	}
	s1 := &Swagger{
		Openapi:    s.Openapi,
		Servers:    s.Servers,
		Info:       s.Info,
		Tags:       s.Tags,
		Paths:      Paths{},
		Components: s.Components,
	}
	for key, path := range s.Paths {
		for _, s2 := range path {
			if strings.ToLower(s2.Resource) == strings.ToLower(resource) {
				s1.Paths[key] = path
				break
			}
		}
	}
	return s1
}

type Server struct {
	Url       string `json:"url,omitempty"`
	Variables struct {
		Scheme string `json:"scheme,omitempty"`
	} `json:"variables,omitempty"`
}

type Info struct {
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Title       string `json:"title,omitempty"`
	Contact     struct {
		Name  string `json:"name,omitempty"`
		Url   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"contact,omitempty"`
}

type Tag struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
type Paths map[string]Path

type Path map[string]*PathMethod

type PathMethod struct {
	Resource    string               `json:"resource,omitempty"`
	Path        string               `json:"-"`
	Method      string               `json:"-"`
	Tags        []string             `json:"tags,omitempty"`
	OperationId string               `json:"operationId,omitempty"`
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`
	Parameters  []*Parameter         `json:"parameters,omitempty"`
}

type Response struct {
	Description string           `json:"description,omitempty"`
	Content     *ResponseContent `json:"content,omitempty"`
}

type ResponseContent struct {
	ApplicationJson *ApplicationJson `json:"application/json,omitempty"`
}

type ApplicationJson struct {
	Schema  *Schema `json:"schema,omitempty"`
	Example any     `json:"example,omitempty"`
}

type SchemaProperties struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Default     string `json:"default,omitempty"`
}

type RequestBody struct {
	Content map[string]*RequestContent `json:"content,omitempty"`
}

type RequestContent struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Components struct {
	Schemas map[string]*Schema `json:"schemas,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name,omitempty"`
	In          string  `json:"in,omitempty"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Example     any     `json:"example,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}
