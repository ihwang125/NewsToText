package repositories

import (
	"news-to-text/internal/models"
	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *models.Alert) error
	GetByID(id uint) (*models.Alert, error)
	GetByUserID(userID uint) ([]models.Alert, error)
	GetActiveAlerts() ([]models.Alert, error)
	Update(alert *models.Alert) error
	Delete(id uint) error
	CreateHistory(history *models.AlertHistory) error
	GetHistoryByAlertID(alertID uint) ([]models.AlertHistory, error)
	GetHistoryByUserID(userID uint) ([]models.AlertHistory, error)
}

type alertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepository) GetByID(id uint) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Preload("User").First(&alert, id).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *alertRepository) GetByUserID(userID uint) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("user_id = ?", userID).Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) GetActiveAlerts() ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("active = ?", true).Preload("User").Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) Update(alert *models.Alert) error {
	return r.db.Save(alert).Error
}

func (r *alertRepository) Delete(id uint) error {
	return r.db.Delete(&models.Alert{}, id).Error
}

func (r *alertRepository) CreateHistory(history *models.AlertHistory) error {
	return r.db.Create(history).Error
}

func (r *alertRepository) GetHistoryByAlertID(alertID uint) ([]models.AlertHistory, error) {
	var history []models.AlertHistory
	err := r.db.Where("alert_id = ?", alertID).Order("created_at DESC").Find(&history).Error
	return history, err
}

func (r *alertRepository) GetHistoryByUserID(userID uint) ([]models.AlertHistory, error) {
	var history []models.AlertHistory
	err := r.db.Joins("JOIN alerts ON alert_histories.alert_id = alerts.id").
		Where("alerts.user_id = ?", userID).
		Order("alert_histories.created_at DESC").
		Find(&history).Error
	return history, err
}