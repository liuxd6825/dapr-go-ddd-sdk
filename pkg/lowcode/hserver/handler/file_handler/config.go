package file_handler

type WebConfig struct {
	Npm *Npm `json:"npm"`
}

type Npm struct {
	Links []*NpmLink `json:"links"`
}

type NpmLink struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
