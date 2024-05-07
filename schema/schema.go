package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Schema struct {
	Schema      string     `json:"$schema,omitempty"`
	Table       string     `json:"table,omitempty"`
	Id          string     `json:"$id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
	Definitions Properties `json:"definitions,omitempty"`
	Required    []string   `json:"required,omitempty"`
}

type Properties map[string]*Property

type Property struct {
	Name             string     `json:"-"`
	Title            string     `json:"title,omitempty"`
	Type             string     `json:"type,omitempty"`
	Properties       Properties `json:"properties,omitempty"`
	Required         []string   `json:"required,omitempty"`
	Description      string     `json:"description,omitempty"`
	Ref              string     `json:"$ref,omitempty"`
	Items            *Items     `json:"items,omitempty"`
	MinItems         *int       `json:"minItems,omitempty"`
	UniqueItems      *bool      `json:"uniqueItems,omitempty"`
	ExclusiveMinimum *int       `json:"exclusiveMinimum,omitempty"`
	Minimum          *int       `json:"minimum,omitempty"`
	Maximum          *int       `json:"maximum,omitempty"`
}

type Items struct {
	Type string `json:"type"`
}

func NewSchemaString(json string) (*Schema, error) {
	var reader = strings.NewReader(json)
	return NewSchema(reader)
}

func NewSchemaFile(pathFile string) (*Schema, error) {
	// 打开 JSON 文件
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	return NewSchema(file)
}

func NewSchema(reader io.Reader) (*Schema, error) {
	var data Schema
	var err error
	// 解码 JSON 数据到结构体中
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&data); err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding JSON: %v", err))
	}
	for key, property := range data.Properties {
		property.Name = key
	}
	return &data, err
}

func (p *Property) GetTitle() string {
	if p.Title == "" {
		return p.Name
	}
	return p.Title
}

func (s *Schema) ToJson() string {
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}
