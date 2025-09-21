-- Initial schema for news-to-text application

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_email (email)
);

-- Alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    topic VARCHAR(255) NOT NULL,
    keywords JSON,
    frequency ENUM('realtime', 'hourly', 'daily') NOT NULL DEFAULT 'daily',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_checked TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_alerts_user_id (user_id),
    INDEX idx_alerts_deleted_at (deleted_at),
    INDEX idx_alerts_active (active),
    INDEX idx_alerts_frequency (frequency),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Alert history table
CREATE TABLE IF NOT EXISTS alert_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    alert_id BIGINT UNSIGNED NOT NULL,
    news_title TEXT NOT NULL,
    news_url TEXT NOT NULL,
    news_source VARCHAR(255),
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    error_msg TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_alert_histories_alert_id (alert_id),
    INDEX idx_alert_histories_sent_at (sent_at),
    INDEX idx_alert_histories_success (success),
    FOREIGN KEY (alert_id) REFERENCES alerts(id) ON DELETE CASCADE
);

-- News sources table
CREATE TABLE IF NOT EXISTS news_sources (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    rss_feed_url TEXT,
    api_endpoint TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_news_sources_active (active),
    INDEX idx_news_sources_category (category),
    INDEX idx_news_sources_deleted_at (deleted_at)
);

-- Insert default news sources
INSERT INTO news_sources (name, url, rss_feed_url, category) VALUES
('TechCrunch', 'https://techcrunch.com', 'https://techcrunch.com/feed/', 'Tech'),
('Reuters Technology', 'https://reuters.com/technology', 'https://feeds.reuters.com/reuters/technologyNews', 'Tech'),
('Bloomberg Markets', 'https://bloomberg.com/markets', 'https://feeds.bloomberg.com/markets/news.rss', 'Stocks'),
('MarketWatch', 'https://marketwatch.com', 'https://feeds.marketwatch.com/marketwatch/marketpulse/', 'Stocks'),
('Hacker News', 'https://news.ycombinator.com', 'https://hnrss.org/frontpage', 'Tech'),
('CNBC Technology', 'https://cnbc.com/technology', 'https://search.cnbc.com/rs/search/combinedcms/view.xml?partnerId=wrss01&id=19854910', 'Tech');