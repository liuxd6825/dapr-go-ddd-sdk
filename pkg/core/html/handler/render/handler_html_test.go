package render

import "testing"

func Test_parserHtml(t *testing.T) {
	html := `
    <script type="module">
        import "ui5e/src/ui5-all.ts"
        import "ui5e/src/index.ts"
        import "ui5e/src/xhtml/index.ts"
        import "ui5e/src/form/form.ts"
    </script>
`
	links := []*NpmLink{}
	links = append(links, &NpmLink{
		Name: "ui5e",
		Path: "/@fs/Users/lxd/Projects/duxm/h-master/packages/ui5e",
	})
	webCfg := &WebConfig{
		Npm: &Npm{
			Links: links,
		},
	}
	h, err := parserHtml(nil, webCfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(h)
}
