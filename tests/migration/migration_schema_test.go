package migration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"gorm.io/gorm"
)

func TestMigrationSchemaContainsPreReleaseFields(t *testing.T) {
	resetDatabaseFlags()
	if os.Getenv("SQL_DSN") == "" {
		sqlitePath := filepath.Join(t.TempDir(), "migration-smoke.db")
		t.Setenv("SQL_DSN", "local")
		common.SQLitePath = sqlitePath
	}
	common.IsMasterNode = true
	common.RedisEnabled = false
	common.BatchUpdateEnabled = false

	if err := model.InitDB(); err != nil {
		t.Fatalf("InitDB migration failed: %v", err)
	}
	if err := model.InitLogDB(); err != nil {
		t.Fatalf("InitLogDB migration failed: %v", err)
	}

	db := model.DB
	logDB := model.LOG_DB

	assertColumnsInDB(t, db, "channels",
		"provider_type",
		"provider_account_id",
		"risk_level",
		"available_scope",
		"visibility",
		"manual_enable_required",
	)
	assertColumnsInDB(t, db, "provider_accounts",
		"id",
		"name",
		"provider_type",
		"encrypted_key",
		"created_at",
	)
	assertColumnsInDB(t, db, "channel_model_mappings",
		"id",
		"channel_id",
		"public_model_name",
		"provider_model_name",
		"enabled",
		"input_price",
		"output_price",
	)
	assertColumnsInDB(t, db, "organizations",
		"id",
		"name",
		"owner_id",
		"status",
	)
	assertColumnsInDB(t, db, "organization_members",
		"id",
		"org_id",
		"user_id",
		"role",
	)
	assertColumnsInDB(t, db, "projects",
		"id",
		"org_id",
		"name",
		"status",
	)
	assertColumnsInDB(t, db, "tokens",
		"org_id",
		"project_id",
		"allow_experimental",
		"allowed_provider_types",
		"key_hash",
		"key_prefix",
	)
	assertColumnsInDB(t, logDB, "logs",
		"org_id",
		"project_id",
		"is_experimental_proxy",
		"provider_type",
		"request_id",
		"upstream_request_id",
	)
}

func resetDatabaseFlags() {
	common.UsingSQLite = false
	common.UsingMySQL = false
	common.UsingPostgreSQL = false
	common.LogSqlType = common.DatabaseTypeSQLite
}

func assertColumnsInDB(t *testing.T, db *gorm.DB, table string, columns ...string) {
	t.Helper()
	if !db.Migrator().HasTable(table) {
		t.Fatalf("expected table %s to exist", table)
	}
	for _, column := range columns {
		if !db.Migrator().HasColumn(table, column) {
			t.Fatalf("expected %s.%s to exist", table, column)
		}
	}
}
