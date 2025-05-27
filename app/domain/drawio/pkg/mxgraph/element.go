package mxgraph

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Element interface {
	GetId() string
	GetXMLName() xml.Name
}

type DrawioFile struct {
	XMLName  xml.Name   `xml:"mxfile"`
	Diagrams []*Diagram `xml:"diagram"`
	Host     string     `xml:"host,attr"`
	Modified string     `xml:"modified,attr"`
	Agent    string     `xml:"agent,attr"`
	ETag     string     `xml:"etag,attr"`
	Version  string     `xml:"version,attr"`
	Pages    int64      `xml:"pages,attr"`
	elements map[string]Element
}

type Diagram struct {
	XMLName    xml.Name   `xml:"diagram"`
	GraphModel GraphModel `xml:"mxGraphModel"`
}

type GraphModel struct {
	XMLName     xml.Name      `xml:"mxGraphModel"`
	Cells       []*MxCell     `xml:"root>mxCell"`
	UserObjects []*UserObject `xml:"root>UserObject"`
	Objects     []*Object     `xml:"root>object"`
}

type MxCell struct {
	XMLName     xml.Name `xml:"mxCell"`
	Id          string   `xml:"id,attr"`
	Edge        string   `xml:"edge,attr"`
	Style       string   `xml:"style,attr"`
	Parent      string   `xml:"parent,attr"`
	Source      string   `xml:"source,attr"`
	Target      string   `xml:"target,attr"`
	Vertex      string   `xml:"vertex,attr"`
	Connectable string   `xml:"connectable,attr"`
}

type Object struct {
	Id       string   `xml:"id,attr"`
	XMLName  xml.Name `xml:"object"`
	Label    string   `xml:"label,attr"`
	CellType string   `xml:"cellType,attr"`
}

type UserObject struct {
	Id      string   `xml:"id,attr"`
	XMLName xml.Name `xml:"UserObject"`
	Label   string   `xml:"label,attr"`
	Link    string   `xml:"link,attr"`
}

func newDrawioFile(file io.Reader) (*DrawioFile, error) {
	decoder := xml.NewDecoder(file)
	var drawio DrawioFile
	err := decoder.Decode(&drawio)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return nil, err
	}
	return &drawio, nil
}

func NewDrawioFile(content string) (*DrawioFile, error) {
	reader := strings.NewReader(content)
	return newDrawioFile(reader)
}

func (d *DrawioFile) GetElement(id string) Element {
	if d.elements == nil {
		d.elements = make(map[string]Element)
		for _, diagrams := range d.Diagrams {
			for _, obj := range diagrams.GraphModel.Objects {
				d.elements[obj.Id] = obj
			}
			for _, userObj := range diagrams.GraphModel.UserObjects {
				d.elements[userObj.Id] = userObj
			}
			for _, cell := range diagrams.GraphModel.Cells {
				d.elements[cell.Id] = cell
			}
		}
	}
	return d.elements[id]
}

func (u *UserObject) GetId() string {
	return u.Id
}

func (u *UserObject) GetXMLName() xml.Name {
	return u.XMLName
}

func (u *MxCell) GetId() string {
	return u.Id
}

func (u *MxCell) GetXMLName() xml.Name {
	return u.XMLName
}

func (u *Object) GetId() string {
	return u.Id
}

func (u *Object) GetXMLName() xml.Name {
	return u.XMLName
}
