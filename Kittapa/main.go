package main

import (
	"context"
)

type (
	Server struct {
		routers map[string]Route
	}

	Route struct {
		Method string      `json:"method"`
		Path   string      `json:"path"`
		Name   HandlerFunc `json:"name"`
	}
)
type HandlerFunc func(context.Context) error

func (s *Server) GET(path string, h HandlerFunc) *Route {
	m := "GET"
	return s.Add(m, path, h)
}

func (s *Server) POST(path string, h HandlerFunc) *Route {
	m := "POST"
	return s.Add(m, path, h)
}

func (s *Server) Add(m, path string, h HandlerFunc) *Route {
	r := Route{
		Method: m,
		Path:   path,
		Name:   h,
	}
	s.routers[path] = r
	return r
}
