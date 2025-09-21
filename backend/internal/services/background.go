package services

import (
	"context"
	"sync"
	"time"

	"news-to-text/internal/models"
	"news-to-text/pkg/logger"
)

type BackgroundService interface {
	Start()
	Stop()
	ProcessAlerts() error
}

type backgroundService struct {
	alertService        AlertService
	newsService         NewsService
	notificationService NotificationService
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	running             bool
	mu                  sync.RWMutex
}

func NewBackgroundService(
	alertService AlertService,
	newsService NewsService,
	notificationService NotificationService,
) BackgroundService {
	ctx, cancel := context.WithCancel(context.Background())

	return &backgroundService{
		alertService:        alertService,
		newsService:         newsService,
		notificationService: notificationService,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

func (s *backgroundService) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	logger.Info("Starting background service...")

	// Start different job processors
	s.wg.Add(3)
	go s.realtimeProcessor()
	go s.hourlyProcessor()
	go s.dailyProcessor()

	s.wg.Wait()
}

func (s *backgroundService) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	logger.Info("Stopping background service...")
	s.cancel()
	s.wg.Wait()
	logger.Info("Background service stopped")
}

func (s *backgroundService) realtimeProcessor() {
	defer s.wg.Done()

	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes for real-time alerts
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			logger.Debug("Processing real-time alerts...")
			if err := s.processAlertsByFrequency(models.FrequencyRealTime); err != nil {
				logger.Error("Error processing real-time alerts:", err)
			}
		}
	}
}

func (s *backgroundService) hourlyProcessor() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			logger.Debug("Processing hourly alerts...")
			if err := s.processAlertsByFrequency(models.FrequencyHourly); err != nil {
				logger.Error("Error processing hourly alerts:", err)
			}
		}
	}
}

func (s *backgroundService) dailyProcessor() {
	defer s.wg.Done()

	// Calculate time until next 9 AM
	now := time.Now()
	next9AM := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if now.After(next9AM) {
		next9AM = next9AM.Add(24 * time.Hour)
	}

	timer := time.NewTimer(time.Until(next9AM))
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			logger.Debug("Processing daily alerts...")
			if err := s.processAlertsByFrequency(models.FrequencyDaily); err != nil {
				logger.Error("Error processing daily alerts:", err)
			}
			// Reset timer for next day
			timer.Reset(24 * time.Hour)
		}
	}
}

func (s *backgroundService) processAlertsByFrequency(frequency models.AlertFrequency) error {
	alerts, err := s.alertService.GetActiveAlerts()
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		if alert.Frequency != frequency {
			continue
		}

		// Check if we should process this alert based on last checked time
		if !s.shouldProcessAlert(&alert, frequency) {
			continue
		}

		if err := s.processAlert(&alert); err != nil {
			logger.Error("Error processing alert", alert.ID, ":", err)
			continue
		}

		// Update last checked time
		if err := s.alertService.UpdateLastChecked(alert.ID); err != nil {
			logger.Error("Error updating last checked time for alert", alert.ID, ":", err)
		}
	}

	return nil
}

func (s *backgroundService) shouldProcessAlert(alert *models.Alert, frequency models.AlertFrequency) bool {
	if alert.LastChecked == nil {
		return true
	}

	now := time.Now()
	timeSinceLastCheck := now.Sub(*alert.LastChecked)

	switch frequency {
	case models.FrequencyRealTime:
		return timeSinceLastCheck >= 5*time.Minute
	case models.FrequencyHourly:
		return timeSinceLastCheck >= 1*time.Hour
	case models.FrequencyDaily:
		return timeSinceLastCheck >= 24*time.Hour
	default:
		return false
	}
}

func (s *backgroundService) processAlert(alert *models.Alert) error {
	logger.Debug("Processing alert:", alert.ID, "Topic:", alert.Topic)

	// Fetch news articles based on keywords
	articles, err := s.newsService.FetchNewsByKeywords(alert.Keywords)
	if err != nil {
		return err
	}

	// Filter articles that were published since last check
	if alert.LastChecked != nil {
		var recentArticles []models.NewsArticle
		for _, article := range articles {
			if article.PublishedAt.After(*alert.LastChecked) {
				recentArticles = append(recentArticles, article)
			}
		}
		articles = recentArticles
	}

	// If no new articles, skip notification
	if len(articles) == 0 {
		logger.Debug("No new articles found for alert:", alert.ID)
		return nil
	}

	logger.Info("Found", len(articles), "new articles for alert:", alert.ID)

	// Send notification
	if err := s.notificationService.SendNewsAlert(&alert.User, alert, articles); err != nil {
		logger.Error("Failed to send notification for alert", alert.ID, ":", err)
		return err
	}

	logger.Info("Successfully sent notification for alert:", alert.ID)
	return nil
}

func (s *backgroundService) ProcessAlerts() error {
	return s.processAlertsByFrequency(models.FrequencyRealTime)
}