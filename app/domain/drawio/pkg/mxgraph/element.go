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

type ElementType string

const (
	ElementType_None       ElementType = ""
	ElementType_MxCell     ElementType = "mxCell"
	ElementType_Object     ElementType = "object"
	ElementType_UserObject ElementType = "UserObject"
)

func (e ElementType) String() string {
	return string(e)
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
	Value       string   `xml:"value,attr"`
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

func (d *DrawioFile) init() {
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
}

func (d *DrawioFile) GetElement(id string) Element {
	d.init()
	return d.elements[id]
}

func (d *DrawioFile) GetMxCellByParent(parent string) []*MxCell {
	d.init()
	var items []*MxCell
	for _, el := range d.elements {
		if cell, ok := el.(*MxCell); ok {
			if cell.Parent == parent {
				items = append(items, cell)
			}
		}
	}
	return items
}

func (d *DrawioFile) GetUserObject(id string) *UserObject {
	el := d.GetElement(id)
	if el == nil {
		return nil
	}
	return AsUserObject(el)
}

func (d *DrawioFile) GetObject(id string) *Object {
	el := d.GetElement(id)
	if el == nil {
		return nil
	}
	return AsObject(el)
}

func (d *DrawioFile) GetMxCell(id string) *MxCell {
	el := d.GetElement(id)
	if el == nil {
		return nil
	}
	return AsMxCell(el)
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

func (u *MxCell) IsEdge() bool {
	return u.Edge == "1"
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

func AsUserObject(el Element) *UserObject {
	if obj, ok := el.(*UserObject); ok {
		return obj
	}
	return nil
}

func AsObject(el Element) *Object {
	if obj, ok := el.(*Object); ok {
		return obj
	}
	return nil
}

func AsMxCell(el Element) *MxCell {
	if obj, ok := el.(*MxCell); ok {
		return obj
	}
	return nil
}
