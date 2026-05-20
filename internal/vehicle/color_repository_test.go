package vehicle

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewVehicleColorRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	if repo == nil {
		t.Fatal("Expected repository to be initialized")
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Red", true, "user1", now, "user2", now).
			AddRow(2, "Blue", false, "user1", now, nil, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color")).
			WillReturnRows(rows)

		colors, err := repo.GetAll()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(colors) != 2 {
			t.Errorf("Expected 2 colors, got %d", len(colors))
		}
		if colors[0].Name != "Red" || colors[1].Name != "Blue" {
			t.Errorf("Unexpected color names: %v", colors)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Red", true, "user1", now, "user2", now).
			AddRow("not-an-int", "Blue", false, "user1", now, nil, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color")).
			WillReturnRows(rows)

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color")).
			WillReturnError(errors.New("query failed"))

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	now := time.Now()

	t.Run("Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Red", true, "user1", now, "user2", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		color, err := repo.GetByID(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if color == nil || color.Name != "Red" {
			t.Errorf("Expected 'Red', got %v", color)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color WHERE id = ?")).
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		color, err := repo.GetByID(99)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if color != nil {
			t.Errorf("Expected nil, got %v", color)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color WHERE id = ?")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetByID(1)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	now := time.Now()
	user := "admin"
	color := &VehicleColor{
		Name:      "Green",
		Status:    true,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_color (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WithArgs(color.Name, color.Status, color.CreatedBy, color.CreatedAt, color.UpdatedBy, color.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(10, 1))

		err := repo.Create(color)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if color.ID != 10 {
			t.Errorf("Expected ID 10, got %d", color.ID)
		}
	})

	t.Run("ExecError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_color (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WillReturnError(errors.New("insert failed"))

		err := repo.Create(color)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("LastInsertIDError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_color (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		err := repo.Create(color)
		if err == nil {
			t.Error("Expected error from LastInsertId, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	now := time.Now()
	updatedBy := "admin"
	updatedAt := now
	color := &VehicleColor{
		ID:        1,
		Name:      "Updated",
		Status:    false,
		UpdatedBy: &updatedBy,
		UpdatedAt: &updatedAt,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE vehicle_color SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?")).
			WithArgs(color.Name, color.Status, color.UpdatedBy, color.UpdatedAt, color.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(color)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vehicle_color WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestGetPaged(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)
	now := time.Now()

	t.Run("SuccessNoSearch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Red", true, "user1", now, "user2", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)

		colors, err := repo.GetPaged(10, 0, "id", "desc")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(colors) != 1 {
			t.Errorf("Expected 1 color, got %d", len(colors))
		}
	})

	t.Run("InvalidSortBy", func(t *testing.T) {
		// sortBy "invalid" should fall back to "id"
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}))

		_, err := repo.GetPaged(10, 0, "invalid", "desc")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("InvalidSortOrder", func(t *testing.T) {
		// sortOrder "invalid" should fall back to "desc"
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}))

		_, err := repo.GetPaged(10, 0, "id", "invalid")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow("not-an-int", "Red", true, "user1", now, "user2", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)

		_, err := repo.GetPaged(10, 0, "id", "desc")
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

		_, err := repo.GetPaged(10, 0, "id", "desc")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleColorRepository(db)

	t.Run("SuccessNoSearch", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM vehicle_color")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		count, err := repo.Count()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if count != 10 {
			t.Errorf("Expected count 10, got %d", count)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("db error"))

		_, err := repo.Count()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
