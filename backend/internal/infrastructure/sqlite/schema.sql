-- テーブル: users
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Initialize database schema

CREATE TABLE IF NOT EXISTS investments (
    id TEXT PRIMARY KEY,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    type TEXT NOT NULL,
    strategy TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS portfolios (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS portfolio_investments (
    portfolio_id TEXT NOT NULL,
    investment_id TEXT NOT NULL,
    PRIMARY KEY (portfolio_id, investment_id),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE,
    FOREIGN KEY (investment_id) REFERENCES investments(id) ON DELETE CASCADE
);

-- イベントストアのテーブル追加
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    event_data TEXT NOT NULL,
    occurred_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_portfolio_user_id ON portfolios(user_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_investments_portfolio_id ON portfolio_investments(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_investments_investment_id ON portfolio_investments(investment_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at ON events(occurred_at);
