package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
	"gorm.io/gorm"
)

type AlertFrequency string

const (
	FrequencyRealTime AlertFrequency = "realtime"
	FrequencyHourly   AlertFrequency = "hourly"
	FrequencyDaily    AlertFrequency = "daily"
)

type Keywords []string

func (k *Keywords) Scan(value interface{}) error {
	if value == nil {
		*k = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, k)
	case string:
		return json.Unmarshal([]byte(v), k)
	}

	return errors.New("cannot scan keywords")
}

func (k Keywords) Value() (driver.Value, error) {
	if k == nil {
		return nil, nil
	}
	return json.Marshal(k)
}

type Alert struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Topic       string         `json:"topic" gorm:"not null"`
	Keywords    Keywords       `json:"keywords" gorm:"type:json"`
	Frequency   AlertFrequency `json:"frequency" gorm:"not null;default:'daily'"`
	Active      bool           `json:"active" gorm:"not null;default:true"`
	LastChecked *time.Time     `json:"last_checked"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AlertHistory []AlertHistory `json:"alert_history,omitempty" gorm:"foreignKey:AlertID"`
}

type AlertHistory struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	AlertID    uint      `json:"alert_id" gorm:"not null;index"`
	NewsTitle  string    `json:"news_title" gorm:"not null"`
	NewsURL    string    `json:"news_url" gorm:"not null"`
	NewsSource string    `json:"news_source"`
	SentAt     time.Time `json:"sent_at"`
	Success    bool      `json:"success" gorm:"not null;default:false"`
	ErrorMsg   string    `json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Alert Alert `json:"alert,omitempty" gorm:"foreignKey:AlertID"`
}

type NewsSource struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	URL         string         `json:"url" gorm:"not null"`
	RSSFeedURL  string         `json:"rss_feed_url"`
	APIEndpoint string         `json:"api_endpoint"`
	Active      bool           `json:"active" gorm:"not null;default:true"`
	Category    string         `json:"category"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type AlertCreateRequest struct {
	Topic     string         `json:"topic" binding:"required"`
	Keywords  []string       `json:"keywords" binding:"required,min=1"`
	Frequency AlertFrequency `json:"frequency" binding:"required,oneof=realtime hourly daily"`
}

type AlertUpdateRequest struct {
	Topic     *string         `json:"topic,omitempty"`
	Keywords  *[]string       `json:"keywords,omitempty"`
	Frequency *AlertFrequency `json:"frequency,omitempty"`
	Active    *bool           `json:"active,omitempty"`
}

type AlertResponse struct {
	ID          uint           `json:"id"`
	Topic       string         `json:"topic"`
	Keywords    []string       `json:"keywords"`
	Frequency   AlertFrequency `json:"frequency"`
	Active      bool           `json:"active"`
	LastChecked *time.Time     `json:"last_checked"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (a *Alert) ToResponse() *AlertResponse {
	return &AlertResponse{
		ID:          a.ID,
		Topic:       a.Topic,
		Keywords:    a.Keywords,
		Frequency:   a.Frequency,
		Active:      a.Active,
		LastChecked: a.LastChecked,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}