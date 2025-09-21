package models

import "time"

type NewsArticle struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	ImageURL    string    `json:"image_url,omitempty"`
	Category    string    `json:"category,omitempty"`
}

// Legacy NewsAPI.org response format (kept for RSS fallback)
type NewsAPIResponse struct {
	Status       string              `json:"status"`
	TotalResults int                 `json:"totalResults"`
	Articles     []NewsAPIArticle    `json:"articles"`
}

type NewsAPIArticle struct {
	Source      NewsAPISource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
	Content     string        `json:"content"`
}

type NewsAPISource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// The News API (thenewsapi.com) response format
type TheNewsAPIResponse struct {
	Meta TheNewsAPIMeta        `json:"meta"`
	Data []TheNewsAPIArticle   `json:"data"`
}

type TheNewsAPIMeta struct {
	Found int    `json:"found"`
	Returned int `json:"returned"`
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
}

type TheNewsAPIArticle struct {
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    string    `json:"keywords"`
	Snippet     string    `json:"snippet"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image_url"`
	Language    string    `json:"language"`
	PublishedAt string    `json:"published_at"`
	Source      string    `json:"source"`
	Categories  []string  `json:"categories"`
	Relevance   int       `json:"relevance_score"`
	Locale      string    `json:"locale"`
}

type RSSFeed struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}