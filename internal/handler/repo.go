package handler

import "app/internal/service"

type UserHandler struct {
	svc service.Service
}

func NewHandler(s service.Service) *UserHandler {
	return &UserHandler{svc: s}
}
