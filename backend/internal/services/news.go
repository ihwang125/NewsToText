package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"news-to-text/internal/models"
	"news-to-text/pkg/logger"
)

type NewsService interface {
	FetchNewsByKeywords(keywords []string) ([]models.NewsArticle, error)
	FetchNewsByCategory(category string) ([]models.NewsArticle, error)
	FetchRSSFeed(url string) ([]models.NewsArticle, error)
	MatchArticles(articles []models.NewsArticle, keywords []string) []models.NewsArticle
}

type newsService struct {
	apiKey string
	client *http.Client
}

func NewNewsService(apiKey string) NewsService {
	return &newsService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *newsService) FetchNewsByKeywords(keywords []string) ([]models.NewsArticle, error) {
	if s.apiKey == "" {
		// Fallback to RSS feeds if no API key
		return s.fetchFromDefaultRSSFeeds(keywords)
	}

	// Use The News API for keyword search
	query := strings.Join(keywords, " ")
	apiURL := fmt.Sprintf("https://api.thenewsapi.com/v1/news/all?api_token=%s&search=%s&limit=50&sort=published_at",
		s.apiKey, url.QueryEscape(query))

	resp, err := s.client.Get(apiURL)
	if err != nil {
		// Fallback to RSS if API fails
		logger.Error("The News API failed, falling back to RSS:", err)
		return s.fetchFromDefaultRSSFeeds(keywords)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback to RSS if API returns error
		logger.Error("The News API returned non-200 status:", resp.StatusCode)
		return s.fetchFromDefaultRSSFeeds(keywords)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.fetchFromDefaultRSSFeeds(keywords)
	}

	var newsResp models.TheNewsAPIResponse
	if err := json.Unmarshal(body, &newsResp); err != nil {
		logger.Error("Failed to parse The News API response:", err)
		return s.fetchFromDefaultRSSFeeds(keywords)
	}

	articles := make([]models.NewsArticle, 0, len(newsResp.Data))
	for _, article := range newsResp.Data {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)

		articles = append(articles, models.NewsArticle{
			Title:       article.Title,
			URL:         article.URL,
			Source:      article.Source,
			Description: article.Description,
			PublishedAt: publishedAt,
			ImageURL:    article.ImageURL,
		})
	}

	return articles, nil
}

func (s *newsService) FetchNewsByCategory(category string) ([]models.NewsArticle, error) {
	if s.apiKey == "" {
		// Fallback to RSS feeds if no API key
		return s.fetchFromDefaultRSSFeeds([]string{category})
	}

	// Map common categories to The News API categories
	categoryMap := map[string]string{
		"technology": "tech",
		"tech":       "tech",
		"business":   "business",
		"sports":     "sports",
		"health":     "health",
		"science":    "science",
		"stocks":     "business", // Map stocks to business category
		"finance":    "business",
	}

	apiCategory, exists := categoryMap[strings.ToLower(category)]
	if !exists {
		apiCategory = "general"
	}

	// Use The News API for category search
	apiURL := fmt.Sprintf("https://api.thenewsapi.com/v1/news/top?api_token=%s&categories=%s&limit=50&sort=published_at",
		s.apiKey, apiCategory)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		// Fallback to RSS if API fails
		logger.Error("The News API failed, falling back to RSS:", err)
		return s.fetchFromDefaultRSSFeeds([]string{category})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback to RSS if API returns error
		logger.Error("The News API returned non-200 status:", resp.StatusCode)
		return s.fetchFromDefaultRSSFeeds([]string{category})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.fetchFromDefaultRSSFeeds([]string{category})
	}

	var newsResp models.TheNewsAPIResponse
	if err := json.Unmarshal(body, &newsResp); err != nil {
		logger.Error("Failed to parse The News API response:", err)
		return s.fetchFromDefaultRSSFeeds([]string{category})
	}

	articles := make([]models.NewsArticle, 0, len(newsResp.Data))
	for _, article := range newsResp.Data {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)

		articles = append(articles, models.NewsArticle{
			Title:       article.Title,
			URL:         article.URL,
			Source:      article.Source,
			Description: article.Description,
			PublishedAt: publishedAt,
			ImageURL:    article.ImageURL,
			Category:    category,
		})
	}

	return articles, nil
}

func (s *newsService) FetchRSSFeed(url string) ([]models.NewsArticle, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed models.RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	articles := make([]models.NewsArticle, 0, len(feed.Items))
	for _, item := range feed.Items {
		publishedAt, _ := time.Parse(time.RFC1123Z, item.PubDate)

		articles = append(articles, models.NewsArticle{
			Title:       item.Title,
			URL:         item.Link,
			Source:      feed.Title,
			Description: item.Description,
			PublishedAt: publishedAt,
		})
	}

	return articles, nil
}

func (s *newsService) MatchArticles(articles []models.NewsArticle, keywords []string) []models.NewsArticle {
	var matched []models.NewsArticle

	for _, article := range articles {
		if s.articleMatches(article, keywords) {
			matched = append(matched, article)
		}
	}

	return matched
}

func (s *newsService) articleMatches(article models.NewsArticle, keywords []string) bool {
	content := strings.ToLower(article.Title + " " + article.Description)

	for _, keyword := range keywords {
		if strings.Contains(content, strings.ToLower(keyword)) {
			return true
		}
	}

	return false
}

func (s *newsService) fetchFromDefaultRSSFeeds(keywords []string) ([]models.NewsArticle, error) {
	rssFeeds := []string{
		"https://techcrunch.com/feed/",
		"https://feeds.reuters.com/reuters/technologyNews",
		"https://hnrss.org/frontpage",
		"https://feeds.bloomberg.com/markets/news.rss",
	}

	var allArticles []models.NewsArticle

	for _, feedURL := range rssFeeds {
		articles, err := s.FetchRSSFeed(feedURL)
		if err != nil {
			logger.Error("Failed to fetch RSS feed:", feedURL, err)
			continue
		}

		// Filter articles by keywords
		matched := s.MatchArticles(articles, keywords)
		allArticles = append(allArticles, matched...)
	}

	return allArticles, nil
}