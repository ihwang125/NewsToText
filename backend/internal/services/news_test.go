package services

import (
	"testing"

	"news-to-text/internal/models"
)

func TestNewsService_MatchArticles(t *testing.T) {
	newsService := NewNewsService("") // Empty API key will use RSS fallback

	articles := []models.NewsArticle{
		{
			Title:       "Apple Releases New iPhone with AI Features",
			Description: "The latest iPhone incorporates machine learning capabilities",
		},
		{
			Title:       "Tesla Stock Surges After Earnings",
			Description: "Electric vehicle company reports strong quarterly results",
		},
		{
			Title:       "Bitcoin Price Reaches New High",
			Description: "Cryptocurrency market sees significant gains",
		},
		{
			Title:       "Google Announces AI Research Breakthrough",
			Description: "New artificial intelligence model shows promising results",
		},
	}

	tests := []struct {
		name            string
		keywords        []string
		expectedMatches int
	}{
		{
			name:            "AI keywords",
			keywords:        []string{"AI", "artificial intelligence"},
			expectedMatches: 2, // Apple iPhone AI and Google AI articles
		},
		{
			name:            "Stock keywords",
			keywords:        []string{"stock", "Tesla"},
			expectedMatches: 1, // Tesla stock article
		},
		{
			name:            "Crypto keywords",
			keywords:        []string{"Bitcoin", "cryptocurrency"},
			expectedMatches: 1, // Bitcoin article
		},
		{
			name:            "No match keywords",
			keywords:        []string{"sports", "NBA"},
			expectedMatches: 0,
		},
		{
			name:            "Case insensitive match",
			keywords:        []string{"APPLE", "tesla"},
			expectedMatches: 2, // Apple and Tesla articles
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := newsService.MatchArticles(articles, tt.keywords)

			if len(matches) != tt.expectedMatches {
				t.Errorf("Expected %d matches but got %d", tt.expectedMatches, len(matches))
				for i, match := range matches {
					t.Logf("Match %d: %s", i+1, match.Title)
				}
			}
		})
	}
}

func TestNewsService_articleMatches(t *testing.T) {
	newsService := newsService{} // Access private method for testing

	article := models.NewsArticle{
		Title:       "Apple Announces New MacBook Pro with M3 Chip",
		Description: "The latest MacBook Pro features advanced machine learning capabilities",
	}

	tests := []struct {
		name     string
		keywords []string
		expected bool
	}{
		{
			name:     "Title match",
			keywords: []string{"Apple"},
			expected: true,
		},
		{
			name:     "Description match",
			keywords: []string{"machine learning"},
			expected: true,
		},
		{
			name:     "Multiple keywords - one match",
			keywords: []string{"Google", "Apple"},
			expected: true,
		},
		{
			name:     "Case insensitive match",
			keywords: []string{"APPLE", "macbook"},
			expected: true,
		},
		{
			name:     "No match",
			keywords: []string{"Tesla", "Bitcoin"},
			expected: false,
		},
		{
			name:     "Partial word match",
			keywords: []string{"Mac"},
			expected: true, // Should match "MacBook"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newsService.articleMatches(article, tt.keywords)

			if result != tt.expected {
				t.Errorf("Expected %v but got %v for keywords %v", tt.expected, result, tt.keywords)
			}
		})
	}
}