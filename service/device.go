package service

import (
	"context"

	"hospital_track/models"
	"hospital_track/repository"
)

type SServiceDevice struct {
	repo repository.IDevice
}

func ServiceDevice(repo repository.IDevice) *SServiceDevice {
	return &SServiceDevice{repo: repo}
}

func (s *SServiceDevice) Registry(ctx context.Context, device *models.SDevice) (id int, err error) {
	if id, err = s.repo.Registry(ctx, device); err != nil {
		return
	}

	if id > 0 {
		if id, err = s.repo.Computed(ctx, id); err != nil {
			return
		}
	}

	return
}
