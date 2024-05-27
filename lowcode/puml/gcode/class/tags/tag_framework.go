package tags

type TagFramework struct {
	BaseTag
	Type TagFrameworkType `json:"type"`
}

type TagFrameworkType string

const (
	None          TagFrameworkType = ""
	EventSourcing TagFrameworkType = "es"
	Repository    TagFrameworkType = "repos"
	MVC           TagFrameworkType = "mvc"
)

func NewTagFramework(text string) (*TagFramework, error) {
	tag := &TagFramework{}
	tag.TagType = TagTypeFramework
	switch text {
	case "es":
		tag.Type = EventSourcing
	case "repos":
		tag.Type = Repository
	case "mvc":
		tag.Type = MVC
	default:
		tag.Type = None
	}
	return tag, nil
}
