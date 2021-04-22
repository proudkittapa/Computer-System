package main

import (
	"context"
)

type (
	Server struct {
		routers map[string]Route
	}

	Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	}
)
type HandlerFunc func(context.Context) error

func (s *Server) GET(path string, h HandlerFunc) string {
	return s.Add(path)
}

func (s *Server) Add(path string) string {

}
