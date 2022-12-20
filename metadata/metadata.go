package metadata

type Metadata interface {
	GetProperties() Properties
	SetProperties(v Properties)
}
