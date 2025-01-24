package server

import (
	"testing"
)

var serverHtml = `
<html>
    <server>
        <services>
            <service url="./service/human/human-service.html"></service>
        </services>
    </server>
</html>
`

func Test_NewServerWidthHTML(t *testing.T) {
	/*
		lfs, err := localfs.NewFs(localfs.Config{Path: "/Users/lxd/Projects/duxm/draw-web/src-server/xsrc", Name: "file"})
		if err != nil {
			t.Fatal(err)
			return
		}
	*/
	server, err := NewServer(nil, []byte(serverHtml), nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(server)
}
