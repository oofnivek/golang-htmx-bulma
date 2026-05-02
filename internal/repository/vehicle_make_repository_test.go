package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"golang-htmx-bulma/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewVehicleMakeRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := NewVehicleMakeRepository(db)
	if repo == nil {
		t.Fatal("Expected repository to be initialized")
	}
}

func TestGetAllMake(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleMakeRepository(db)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Toyota", true, "user1", now, "user1", now).
			AddRow(2, "Honda", false, "user1", now, "user1", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make")).
			WillReturnRows(rows)

		makes, err := repo.GetAll()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(makes) != 2 {
			t.Errorf("Expected 2 makes, got %d", len(makes))
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow("not-an-int", "Toyota", true, "user1", now, "user1", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make")).
			WillReturnRows(rows)

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make")).
			WillReturnError(errors.New("query failed"))

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestGetByIDMake(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleMakeRepository(db)
	now := time.Now()

	t.Run("Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "status", "created_by", "created_at", "updated_by", "updated_at"}).
			AddRow(1, "Toyota", true, "user1", now, "user1", now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		make, err := repo.GetByID(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if make == nil || make.Name != "Toyota" {
			t.Errorf("Expected 'Toyota', got %v", make)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make WHERE id = ?")).
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		make, err := repo.GetByID(99)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if make != nil {
			t.Errorf("Expected nil, got %v", make)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make WHERE id = ?")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetByID(1)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreateMakeRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleMakeRepository(db)
	now := time.Now()
	user := "admin"
	make := &model.VehicleMake{
		Name:      "Mazda",
		Status:    true,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_make (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WithArgs(make.Name, make.Status, make.CreatedBy, make.CreatedAt, make.UpdatedBy, make.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(10, 1))

		err := repo.Create(make)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if make.ID != 10 {
			t.Errorf("Expected ID 10, got %d", make.ID)
		}
	})

	t.Run("ExecError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_make (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WillReturnError(errors.New("insert failed"))

		err := repo.Create(make)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("LastInsertIDError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vehicle_make (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)")).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		err := repo.Create(make)
		if err == nil {
			t.Error("Expected error from LastInsertId, got nil")
		}
	})
}

func TestUpdateMakeRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleMakeRepository(db)
	now := time.Now()
	updatedBy := "admin"
	updatedAt := now
	make := &model.VehicleMake{
		ID:        1,
		Name:      "Updated",
		Status:    false,
		UpdatedBy: &updatedBy,
		UpdatedAt: &updatedAt,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE vehicle_make SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?")).
			WithArgs(make.Name, make.Status, make.UpdatedBy, make.UpdatedAt, make.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(make)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestDeleteMakeRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleMakeRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vehicle_make WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
