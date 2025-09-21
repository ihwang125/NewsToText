package services

import (
	"errors"
	"time"

	"news-to-text/internal/models"
	"news-to-text/internal/repositories"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AlertService interface {
	CreateAlert(userID uint, req *models.AlertCreateRequest) (*models.AlertResponse, error)
	GetAlerts(userID uint) ([]models.AlertResponse, error)
	GetAlertByID(userID uint, alertID uint) (*models.AlertResponse, error)
	UpdateAlert(userID uint, alertID uint, req *models.AlertUpdateRequest) (*models.AlertResponse, error)
	DeleteAlert(userID uint, alertID uint) error
	GetAlertHistory(userID uint) ([]models.AlertHistory, error)
	TestAlert(userID uint, alertID uint) error
	GetActiveAlerts() ([]models.Alert, error)
	UpdateLastChecked(alertID uint) error
}

type alertService struct {
	alertRepo repositories.AlertRepository
	redis     *redis.Client
}

func NewAlertService(alertRepo repositories.AlertRepository, redisClient *redis.Client) AlertService {
	return &alertService{
		alertRepo: alertRepo,
		redis:     redisClient,
	}
}

func (s *alertService) CreateAlert(userID uint, req *models.AlertCreateRequest) (*models.AlertResponse, error) {
	alert := &models.Alert{
		UserID:    userID,
		Topic:     req.Topic,
		Keywords:  models.Keywords(req.Keywords),
		Frequency: req.Frequency,
		Active:    true,
	}

	if err := s.alertRepo.Create(alert); err != nil {
		return nil, err
	}

	return alert.ToResponse(), nil
}

func (s *alertService) GetAlerts(userID uint) ([]models.AlertResponse, error) {
	alerts, err := s.alertRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.AlertResponse, len(alerts))
	for i, alert := range alerts {
		responses[i] = *alert.ToResponse()
	}

	return responses, nil
}

func (s *alertService) GetAlertByID(userID uint, alertID uint) (*models.AlertResponse, error) {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alert not found")
		}
		return nil, err
	}

	if alert.UserID != userID {
		return nil, errors.New("unauthorized access to alert")
	}

	return alert.ToResponse(), nil
}

func (s *alertService) UpdateAlert(userID uint, alertID uint, req *models.AlertUpdateRequest) (*models.AlertResponse, error) {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alert not found")
		}
		return nil, err
	}

	if alert.UserID != userID {
		return nil, errors.New("unauthorized access to alert")
	}

	// Update fields
	if req.Topic != nil {
		alert.Topic = *req.Topic
	}
	if req.Keywords != nil {
		alert.Keywords = models.Keywords(*req.Keywords)
	}
	if req.Frequency != nil {
		alert.Frequency = *req.Frequency
	}
	if req.Active != nil {
		alert.Active = *req.Active
	}

	if err := s.alertRepo.Update(alert); err != nil {
		return nil, err
	}

	return alert.ToResponse(), nil
}

func (s *alertService) DeleteAlert(userID uint, alertID uint) error {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("alert not found")
		}
		return err
	}

	if alert.UserID != userID {
		return errors.New("unauthorized access to alert")
	}

	return s.alertRepo.Delete(alertID)
}

func (s *alertService) GetAlertHistory(userID uint) ([]models.AlertHistory, error) {
	return s.alertRepo.GetHistoryByUserID(userID)
}

func (s *alertService) TestAlert(userID uint, alertID uint) error {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		return err
	}

	if alert.UserID != userID {
		return errors.New("unauthorized access to alert")
	}

	// Create a test history entry
	history := &models.AlertHistory{
		AlertID:    alertID,
		NewsTitle:  "Test Alert - " + alert.Topic,
		NewsURL:    "https://example.com/test-news",
		NewsSource: "Test Source",
		Success:    true,
		SentAt:     time.Now(),
	}

	return s.alertRepo.CreateHistory(history)
}

func (s *alertService) GetActiveAlerts() ([]models.Alert, error) {
	return s.alertRepo.GetActiveAlerts()
}

func (s *alertService) UpdateLastChecked(alertID uint) error {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		return err
	}

	now := time.Now()
	alert.LastChecked = &now
	return s.alertRepo.Update(alert)
}