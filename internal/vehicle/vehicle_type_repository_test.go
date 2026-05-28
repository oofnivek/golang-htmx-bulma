package vehicle

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var typeColumns = []string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}

func TestNewVehicleTypeRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	if repo == nil {
		t.Fatal("expected repository to be initialized")
	}
}

func TestVehicleTypeGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows(typeColumns).
			AddRow(1, "Sedan", true, "user1", now, "user2", now).
			AddRow(2, "SUV", false, "user1", now, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id DESC")).
			WillReturnRows(rows)

		types, err := repo.GetAll()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(types) != 2 {
			t.Errorf("expected 2, got %d", len(types))
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id DESC")).
			WillReturnError(errors.New("query failed"))
		_, err := repo.GetAll()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows(typeColumns).
			AddRow("not-an-int", "Sedan", true, "user1", now, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id DESC")).
			WillReturnRows(rows)
		_, err := repo.GetAll()
		if err == nil {
			t.Error("expected scan error, got nil")
		}
	})
}

func TestVehicleTypeGetPaged(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows(typeColumns).
			AddRow(1, "Sedan", true, "user1", now, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)
		types, err := repo.GetPaged(10, 0, "id", "desc")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(types) != 1 {
			t.Errorf("expected 1, got %d", len(types))
		}
	})

	t.Run("InvalidSortBy", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows(typeColumns))
		_, err := repo.GetPaged(10, 0, "invalid", "desc")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("InvalidSortOrder", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows(typeColumns))
		_, err := repo.GetPaged(10, 0, "id", "invalid")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
		_, err := repo.GetPaged(10, 0, "id", "desc")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows(typeColumns).
			AddRow("not-an-int", "Sedan", true, "user1", now, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type ORDER BY id desc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)
		_, err := repo.GetPaged(10, 0, "id", "desc")
		if err == nil {
			t.Error("expected scan error, got nil")
		}
	})
}

func TestVehicleTypeCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM vehicle_type")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
		count, err := repo.Count()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 5 {
			t.Errorf("expected 5, got %d", count)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM vehicle_type")).
			WillReturnError(errors.New("db error"))
		_, err := repo.Count()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestVehicleTypeGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	now := time.Now()

	t.Run("Found", func(t *testing.T) {
		rows := sqlmock.NewRows(typeColumns).
			AddRow(1, "Sedan", true, "user1", now, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnRows(rows)
		vt, err := repo.GetByID(1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if vt == nil || vt.Name != "Sedan" {
			t.Errorf("expected Sedan, got %v", vt)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type WHERE id = ?")).
			WithArgs(int64(99)).
			WillReturnError(sql.ErrNoRows)
		vt, err := repo.GetByID(99)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if vt != nil {
			t.Errorf("expected nil, got %v", vt)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_type WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db error"))
		_, err := repo.GetByID(1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestVehicleTypeCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	now := time.Now()
	user := "admin"
	vt := &VehicleType{
		Name:      "Sedan",
		Status:    true,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_type (name, status, created_by, created_at, updated_by, updated_at) VALUES (?,?,?,?,?,?)")).
			WithArgs(vt.Name, vt.Status, vt.CreatedBy, vt.CreatedAt, vt.UpdatedBy, vt.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(42, 1))
		err := repo.Create(vt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if vt.ID != 42 {
			t.Errorf("expected ID 42, got %d", vt.ID)
		}
	})

	t.Run("ExecError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_type (name, status, created_by, created_at, updated_by, updated_at) VALUES (?,?,?,?,?,?)")).
			WillReturnError(errors.New("insert failed"))
		err := repo.Create(vt)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("LastInsertIDError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_type (name, status, created_by, created_at, updated_by, updated_at) VALUES (?,?,?,?,?,?)")).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
		err := repo.Create(vt)
		if err == nil {
			t.Error("expected error from LastInsertId, got nil")
		}
	})
}

func TestVehicleTypeUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)
	now := time.Now()
	user := "admin"
	vt := &VehicleType{
		ID:        1,
		Name:      "Updated",
		Status:    false,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE vehicle_type SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?")).
			WithArgs(vt.Name, vt.Status, vt.UpdatedBy, vt.UpdatedAt, vt.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.Update(vt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE vehicle_type SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?")).
			WillReturnError(errors.New("update failed"))
		err := repo.Update(vt)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestVehicleTypeDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()
	repo := NewVehicleTypeRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vehicle_type WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		err := repo.Delete(1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vehicle_type WHERE id = ?")).
			WillReturnError(errors.New("delete failed"))
		err := repo.Delete(1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
