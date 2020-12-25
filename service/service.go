package service

import (
	"context"
	"hospital_track/models"
	"hospital_track/repository"
)

// IDevice интерфейс с методами устройства
type IDevice interface {
	Registry(ctx context.Context, dev *models.SDevice) (int, error)
}

// SService структура сервиса
type SService struct {
	IDevice
}

// Service конструктор сервиса
func Service(repos *repository.SRepository) *SService {
	return &SService{
		IDevice: ServiceDevice(repos.IDevice),
	}

}
