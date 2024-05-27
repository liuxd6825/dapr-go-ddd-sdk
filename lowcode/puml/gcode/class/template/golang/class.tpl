package {{.Package}}
{{ $className:=.Name }}
{{ $sn :=.SN }}
type {{.Name}} struct {
    {{- range $index, $field := .Fields }}
    {{$field.Name}} {{$field.DataType}}
    {{- end}}
}

//
// New{{.Name}}
// @Description: {{.NotesText}}
func New{{.Name}}() *{{.Name}} {
    return &{{.Name}}{}
}

{{range $method := .Methods }}

//
// {{if $method.IsPublic}}{{$method.FUpperName}}{{else}}{{$method.FLowerName}}{{end}}
// @Description:{{$method.NotesText}}
{{- range $arg := $method.Args }}
// @param {{$arg.Name}}
{{- end}}
{{- range $res := $method.Results }}
// @return {{$res.Name}} {{$res.DataType}}
{{- end}}
func ({{$sn}} *{{$className}}) {{if $method.IsPublic}}{{$method.FUpperName}}{{else}}{{$method.FLowerName}}{{end}}({{range $idx,$arg:=$method.Args}}{{$arg.Name}} {{$arg.DataType}}{{if $method.NotEndArg $idx }},{{end}}{{end}}){{if $method.HasResults}}{{if gt $method.LenResults 1 }}({{end}}{{range $idx,$res:=$method.Results}}{{$res.Name}} {{$res.DataType}}{{if $method.NotEndResult $idx}}, {{end}}{{end}}{{if gt $method.LenResults 1}}){{end}}{{end}} {
}
{{end}}