package handlers

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	healthzv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/healthz/v1"
	"go.breu.io/quantm/internal/nomad/proto/ctrlplane/healthz/v1/healthzv1connect"
)

type (
	HealthCheckService struct{}
)

func (s *HealthCheckService) Status(
	ctx context.Context, _ *connect.Request[emptypb.Empty],
) (*connect.Response[healthzv1.StatusResponse], error) {
	response := connect.NewResponse(&healthzv1.StatusResponse{
		Database: true,
		Temporal: true,
	})

	return response, nil
}

func NewHealthCheckServiceHandler() (string, http.Handler) {
	return healthzv1connect.NewHealthCheckServiceHandler(&HealthCheckService{})
}
