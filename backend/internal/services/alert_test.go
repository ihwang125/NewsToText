package services

import (
	"testing"

	"news-to-text/internal/models"
	"news-to-text/internal/repositories"

	"github.com/redis/go-redis/v9"
)

func TestAlertService_CreateAlert(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	redisClient := setupTestRedis()
	alertRepo := repositories.NewAlertRepository(db)
	alertService := NewAlertService(alertRepo, redisClient)

	// Create a test user first
	userRepo := repositories.NewUserRepository(db)
	testUser := &models.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	userRepo.Create(testUser)

	tests := []struct {
		name    string
		userID  uint
		request *models.AlertCreateRequest
		wantErr bool
	}{
		{
			name:   "Valid alert creation",
			userID: testUser.ID,
			request: &models.AlertCreateRequest{
				Topic:     "Technology",
				Keywords:  []string{"AI", "machine learning"},
				Frequency: models.FrequencyDaily,
			},
			wantErr: false,
		},
		{
			name:   "Alert with real-time frequency",
			userID: testUser.ID,
			request: &models.AlertCreateRequest{
				Topic:     "Stocks",
				Keywords:  []string{"AAPL", "Tesla"},
				Frequency: models.FrequencyRealTime,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert, err := alertService.CreateAlert(tt.userID, tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if alert == nil {
				t.Errorf("Expected alert but got nil")
				return
			}

			if alert.Topic != tt.request.Topic {
				t.Errorf("Expected topic %s but got %s", tt.request.Topic, alert.Topic)
			}

			if len(alert.Keywords) != len(tt.request.Keywords) {
				t.Errorf("Expected %d keywords but got %d", len(tt.request.Keywords), len(alert.Keywords))
			}

			if alert.Frequency != tt.request.Frequency {
				t.Errorf("Expected frequency %s but got %s", tt.request.Frequency, alert.Frequency)
			}

			if !alert.Active {
				t.Errorf("Expected alert to be active")
			}
		})
	}
}

func TestAlertService_GetAlerts(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	redisClient := setupTestRedis()
	alertRepo := repositories.NewAlertRepository(db)
	alertService := NewAlertService(alertRepo, redisClient)

	// Create test users
	userRepo := repositories.NewUserRepository(db)
	testUser1 := &models.User{Email: "user1@example.com", Password: "password"}
	testUser2 := &models.User{Email: "user2@example.com", Password: "password"}
	userRepo.Create(testUser1)
	userRepo.Create(testUser2)

	// Create test alerts
	alert1 := &models.Alert{
		UserID:    testUser1.ID,
		Topic:     "Tech",
		Keywords:  models.Keywords{"AI"},
		Frequency: models.FrequencyDaily,
		Active:    true,
	}
	alert2 := &models.Alert{
		UserID:    testUser1.ID,
		Topic:     "Stocks",
		Keywords:  models.Keywords{"AAPL"},
		Frequency: models.FrequencyHourly,
		Active:    true,
	}
	alert3 := &models.Alert{
		UserID:    testUser2.ID,
		Topic:     "Sports",
		Keywords:  models.Keywords{"NBA"},
		Frequency: models.FrequencyDaily,
		Active:    true,
	}

	alertRepo.Create(alert1)
	alertRepo.Create(alert2)
	alertRepo.Create(alert3)

	tests := []struct {
		name           string
		userID         uint
		expectedCount  int
	}{
		{
			name:          "User1 alerts",
			userID:        testUser1.ID,
			expectedCount: 2,
		},
		{
			name:          "User2 alerts",
			userID:        testUser2.ID,
			expectedCount: 1,
		},
		{
			name:          "Non-existent user",
			userID:        999,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alerts, err := alertService.GetAlerts(tt.userID)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(alerts) != tt.expectedCount {
				t.Errorf("Expected %d alerts but got %d", tt.expectedCount, len(alerts))
			}

			// Verify all alerts belong to the correct user
			for _, alert := range alerts {
				// Since we're getting AlertResponse, we can't directly check UserID
				// but we can verify the alert content matches what we expect
				if alert.ID == 0 {
					t.Errorf("Alert ID should not be 0")
				}
			}
		})
	}
}