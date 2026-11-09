package service

import (
	"context"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

type WorkCenterRepository interface {
	List(ctx context.Context) ([]model.WorkCenter, error)
	Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error)
}
type WorkCenter struct {
	repo WorkCenterRepository
}

func NewWorkCenter(repo WorkCenterRepository) *WorkCenter {
	return &WorkCenter{repo: repo}
}

func (s *WorkCenter) List(ctx context.Context) ([]model.WorkCenter, error) {
	return s.repo.List(ctx)
}

func (s *WorkCenter) Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error) {
	req.OperationTypes = dedup(req.OperationTypes)
	return s.repo.Create(ctx, req)
}

func dedup(in []string) []string {
	if len(in) < 2 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
