package sqlite

import (
	"database/sql"
	"moneyget/internal/domain"
	"os"
	"testing"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	dbPath := "test.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// スキーマの作成
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

func TestEventStore_Store(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := NewEventStoreDB(db)

	// テストケース1: InvestmentCreatedEventの保存
	money, err := domain.NewMoney(1000, "JPY")
	if err != nil {
		t.Fatalf("Failed to create money: %v", err)
	}

	investmentID := domain.NewInvestmentID("test-id")
	event := domain.NewInvestmentCreatedEvent(investmentID, money)

	err = store.Store(event)
	if err != nil {
		t.Errorf("Failed to store event: %v", err)
	}

	// 保存されたイベントの検証
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE event_type = ?", "InvestmentCreated").Scan(&count)
	if err != nil {
		t.Errorf("Failed to query events: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 event, got %d", count)
	}

	// テストケース2: PortfolioUpdatedEventの保存
	portfolioID := domain.NewPortfolioID("portfolio-id")
	portfolioEvent := domain.NewPortfolioUpdatedEvent(portfolioID, money)

	err = store.Store(portfolioEvent)
	if err != nil {
		t.Errorf("Failed to store portfolio event: %v", err)
	}

	// 保存されたイベントの検証
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE event_type = ?", "PortfolioUpdated").Scan(&count)
	if err != nil {
		t.Errorf("Failed to query events: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 event, got %d", count)
	}
}
