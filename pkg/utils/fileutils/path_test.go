package fileutils

import "testing"

func Test_absFile(t *testing.T) {

	file := AbsPath("/home/user/", "../../test.html")
	if file != "/test.html" {
		t.Errorf("absFile failed, expect ../../test.html, got %s", file)
	}

	file = AbsPath("/home/user/", "../test.html")
	if file != "/home/test.html" {
		t.Errorf("absFile failed, expect ../test.html, got %s", file)
	}

	file = AbsPath("/home/user/", "./test.html")
	if file != "/home/user/test.html" {
		t.Errorf("absFile failed, expect ./test.html, got %s", file)
	}
}
