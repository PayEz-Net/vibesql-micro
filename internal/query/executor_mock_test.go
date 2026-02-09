package query

import (
	"database/sql"
	"testing"
)

func TestNewExecutor_NotNil(t *testing.T) {
	db := &sql.DB{}
	
	executor := NewExecutor(db)
	
	if executor == nil {
		t.Fatal("NewExecutor returned nil")
	}
	
	if executor.db != db {
		t.Error("executor db field not set correctly")
	}
}

func TestNewExecutor_NilDB(t *testing.T) {
	executor := NewExecutor(nil)
	
	if executor == nil {
		t.Fatal("NewExecutor should not return nil even with nil db")
	}
	
	if executor.db != nil {
		t.Error("executor db should be nil when passed nil")
	}
}
