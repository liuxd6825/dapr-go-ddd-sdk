package values

import "fmt"

var console = make(map[string]any)

func Console() any {
	return console
}

func init() {
	console["log"] = log
	console["info"] = info
	console["warn"] = warn
	console["error"] = error
	console["debug"] = debug
	console["trace"] = trace
}

func log(a ...any) {
	fmt.Println(a...)
}

func info(a ...any) {
	fmt.Println(a...)
}

func warn(a ...any) {
	fmt.Println(a...)
}

func error(a ...any) {
	fmt.Println(a...)
}

func debug(a ...any) {
	fmt.Println(a...)
}

func trace(a ...any) {
	fmt.Println(a...)
}
