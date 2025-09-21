package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"news-to-text/internal/models"
	"news-to-text/pkg/logger"
)

type NotificationService interface {
	SendSMS(phoneNumber, message string) error
	SendNewsAlert(user *models.User, alert *models.Alert, articles []models.NewsArticle) error
	FormatNewsMessage(alert *models.Alert, articles []models.NewsArticle) string
}

type notificationService struct {
	smsAPIKey string
	client    *http.Client
}

func NewNotificationService(smsAPIKey string) NotificationService {
	return &notificationService{
		smsAPIKey: smsAPIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *notificationService) SendSMS(phoneNumber, message string) error {
	if s.smsAPIKey == "" {
		// Mock SMS sending for development
		logger.Info("Mock SMS sent to", phoneNumber, ":", message)
		return nil
	}

	// Example implementation for Twilio
	// You can replace this with your preferred SMS provider
	payload := map[string]string{
		"To":   phoneNumber,
		"From": "+1234567890", // Your Twilio phone number
		"Body": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/YOUR_ACCOUNT_SID/Messages.json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("YOUR_ACCOUNT_SID", s.smsAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *notificationService) SendNewsAlert(user *models.User, alert *models.Alert, articles []models.NewsArticle) error {
	message := s.FormatNewsMessage(alert, articles)

	// For now, we'll assume users have phone numbers in their email field
	// In a real implementation, you'd have a separate phone field
	phoneNumber := "+1234567890" // Placeholder

	return s.SendSMS(phoneNumber, message)
}

func (s *notificationService) FormatNewsMessage(alert *models.Alert, articles []models.NewsArticle) string {
	if len(articles) == 0 {
		return fmt.Sprintf("No new articles found for your alert: %s", alert.Topic)
	}

	message := fmt.Sprintf("🔔 News Alert: %s\n\n", alert.Topic)

	// Limit to top 3 articles to keep message short
	maxArticles := 3
	if len(articles) > maxArticles {
		articles = articles[:maxArticles]
	}

	for i, article := range articles {
		message += fmt.Sprintf("%d. %s\n%s\n\n", i+1, article.Title, article.URL)
	}

	if len(articles) == maxArticles {
		message += "... and more"
	}

	return message
}