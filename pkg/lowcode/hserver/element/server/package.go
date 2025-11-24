package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/ctx_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/db_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/events_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/feign_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/html_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/json_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/params_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/strings_pkg"

	"strings"
)

type PkgSetup interface {
	Add(name string, fun func() any)
	Setup(keys ...string)
}

type pkgSetup struct {
	server *Server
	keys   map[string]func() any
}

func NewPkgSetup(server *Server) PkgSetup {
	s := &pkgSetup{server: server, keys: make(map[string]func() any)}
	s.Add("db", s.db)
	s.Add("feign", s.feign)
	s.Add("fs", s.fs)
	s.Add("json", s.json)
	s.Add("params", s.params)
	s.Add("context", s.context)
	s.Add("schema", s.schema)
	s.Add("strings", s.strings)
	s.Add("env", s.env)
	s.Add("html", s.html)
	s.Add("logs", s.logs)
	return s
}

func (p *pkgSetup) Add(name string, fun func() any) {
	if _, ok := p.keys[name]; !ok {
		p.keys[name] = fun
	}
}

func (p *pkgSetup) Setup(keys ...string) {
	for _, key := range keys {
		key = strings.ToLower(key)
		if !p.server.Pkg().Has(key) {
			if newFun, ok := p.keys[key]; ok && newFun != nil {
				p.server.AddPkg(key, newFun())
			}
		}

		if key == "all" {
			for k, v := range p.keys {
				p.server.AddPkg(k, v())
			}
			break
		}
	}
}

func (p *pkgSetup) db() any {
	return db_pkg.New(p.server)
}

func (p *pkgSetup) feign() any {
	return feign_pkg.New(p.server)
}

func (p *pkgSetup) fs() any {
	return p.server.FsPkg()
}

func (p *pkgSetup) context() any {
	return ctx_pkg.New()
}

func (p *pkgSetup) schema() any {
	return schema_pkg.New(p.server)
}

func (p *pkgSetup) params() any {
	return params_pkg.New(p.server)
}

func (p *pkgSetup) env() any {
	return p.server.env
}

func (p *pkgSetup) logs() any {
	return p.server.Logger()
}

func (p *pkgSetup) json() any {
	return json_pkg.New()
}

func (p *pkgSetup) strings() any {
	return strings_pkg.New()
}

func (p *pkgSetup) html() any {
	return html_pkg.New(p.server)
}

func (p *pkgSetup) events() any {
	return events_pkg.New(p.server)
}
