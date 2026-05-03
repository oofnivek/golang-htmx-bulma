package repository

import (
	"errors"
	"regexp"
	"testing"

	"golang-htmx-bulma/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewRoleRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := NewRoleRepository(db)
	if repo == nil {
		t.Fatal("Expected repository to be initialized")
	}
}

func TestGetAllRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("ADMIN", "Administrator").
			AddRow("USER", "User")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id ASC")).
			WillReturnRows(rows)

		roles, err := repo.GetAll()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(roles) != 2 {
			t.Errorf("Expected 2 roles, got %d", len(roles))
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("OnlyOneColumn")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id ASC")).
			WillReturnRows(rows)

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id ASC")).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetAll()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestGetPagedRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)

	t.Run("SuccessWithoutSearch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("ADMIN", "Administrator")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id asc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)

		roles, err := repo.GetPaged(10, 0, "id", "asc", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
	})

	t.Run("SuccessWithSearch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("ADMIN", "Administrator")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE) ORDER BY name desc LIMIT ? OFFSET ?")).
			WithArgs("admin", 5, 0).
			WillReturnRows(rows)

		roles, err := repo.GetPaged(5, 0, "name", "desc", "admin")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
	})

	t.Run("InvalidSort", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id asc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		_, err := repo.GetPaged(10, 0, "invalid", "invalid", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("fail"))
		_, err := repo.GetPaged(10, 0, "id", "asc", "")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("OnlyOne")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles ORDER BY id asc LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)

		_, err := repo.GetPaged(10, 0, "id", "asc", "")
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})
}

func TestCountRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)

	t.Run("SuccessNoSearch", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM roles")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		count, err := repo.Count("")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	})

	t.Run("SuccessWithSearch", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM roles WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)")).
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		count, err := repo.Count("admin")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected count 1, got %d", count)
		}
	})
}

func TestGetByIDRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)

	t.Run("Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("ADMIN", "Administrator")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM roles WHERE id = ?")).
			WithArgs("ADMIN").
			WillReturnRows(rows)

		role, err := repo.GetByID("ADMIN")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if role == nil || role.ID != "ADMIN" {
			t.Errorf("Expected ADMIN, got %v", role)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("NONE").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
		role, err := repo.GetByID("NONE")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if role != nil {
			t.Error("Expected nil, got role")
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("FAIL").WillReturnError(errors.New("db error"))
		_, err := repo.GetByID("FAIL")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreateRoleRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)
	role := &model.Role{ID: "TEST", Name: "Test Role"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO roles (id, name) VALUES (?, ?)")).
		WithArgs(role.ID, role.Name).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Create(role)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUpdateRoleRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)
	role := &model.Role{ID: "ADMIN", Name: "New Admin"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE roles SET name = ? WHERE id = ?")).
		WithArgs(role.Name, role.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(role)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDeleteRoleRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	repo := NewRoleRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM roles WHERE id = ?")).
		WithArgs("ADMIN").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete("ADMIN")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
