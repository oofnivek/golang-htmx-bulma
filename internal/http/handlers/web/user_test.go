package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-htmx-bulma/internal/user"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func setupUserTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.HTMLRender = view.NewRenderer("../../../../templates")
	return r
}

func TestListUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := &MockUserService{
			ListPagedFn: func(page, pageSize int, sortBy, sortOrder string) ([]user.User, int, error) {
				return []user.User{{Email: "a@b.com", FirstName: "Alice"}}, 1, nil
			},
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users", h.List)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			ListPagedFn: func(page, pageSize int, sortBy, sortOrder string) ([]user.User, int, error) {
				return nil, 0, errors.New("fail")
			},
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users", h.List)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestCreateFormUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := &MockUserService{}
		roleSvc := &MockRoleService{
			ListAllFn: func() ([]user.Role, error) {
				return []user.Role{{ID: "ADMIN", Name: "Administrator"}}, nil
			},
		}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/new", h.CreateForm)

		req, _ := http.NewRequest("GET", "/users/new", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("RoleError", func(t *testing.T) {
		userSvc := &MockUserService{}
		roleSvc := &MockRoleService{
			ListAllFn: func() ([]user.Role, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/new", h.CreateForm)

		req, _ := http.NewRequest("GET", "/users/new", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := &MockUserService{
			CreateUserFn: func(u *user.User) error { return nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.POST("/users", h.Create)

		req, _ := http.NewRequest("POST", "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			CreateUserFn: func(u *user.User) error { return errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.POST("/users", h.Create)

		req, _ := http.NewRequest("POST", "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestEditFormUser(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) {
				return &user.User{Email: email, FirstName: "Alice"}, nil
			},
		}
		roleSvc := &MockRoleService{
			ListAllFn: func() ([]user.Role, error) {
				return []user.Role{{ID: "ADMIN", Name: "Administrator"}}, nil
			},
		}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/users/alice@example.com/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/users/nobody@example.com/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("UserError", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/users/fail@example.com/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("RoleError", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) {
				return &user.User{Email: email}, nil
			},
		}
		roleSvc := &MockRoleService{
			ListAllFn: func() ([]user.Role, error) { return nil, errors.New("fail") },
		}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/users/alice@example.com/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestPasswordFieldsUser(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}
	h := NewUserHandler(userSvc, roleSvc)
	r := setupUserTestRouter()
	r.GET("/users/password-fields", h.PasswordFields)

	req, _ := http.NewRequest("GET", "/users/password-fields", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestViewUser(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) {
				return &user.User{Email: email, FirstName: "Alice"}, nil
			},
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email", h.View)

		req, _ := http.NewRequest("GET", "/users/alice@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email", h.View)

		req, _ := http.NewRequest("GET", "/users/nobody@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email", h.View)

		req, _ := http.NewRequest("GET", "/users/fail@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := &MockUserService{
			UpdateUserFn: func(u *user.User) error { return nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.POST("/users/:email", h.Update)

		req, _ := http.NewRequest("POST", "/users/alice@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			UpdateUserFn: func(u *user.User) error { return errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.POST("/users/:email", h.Update)

		req, _ := http.NewRequest("POST", "/users/alice@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := &MockUserService{
			DeleteUserFn: func(email string) error { return nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.DELETE("/users/:email", h.Delete)

		req, _ := http.NewRequest("DELETE", "/users/alice@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			DeleteUserFn: func(email string) error { return errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.DELETE("/users/:email", h.Delete)

		req, _ := http.NewRequest("DELETE", "/users/alice@example.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteConfirmUser(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) {
				return &user.User{Email: email, FirstName: "Alice", LastName: "Smith"}, nil
			},
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/users/alice@example.com/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, nil },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/users/nobody@example.com/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userSvc := &MockUserService{
			GetByEmailFn: func(email string) (*user.User, error) { return nil, errors.New("fail") },
		}
		roleSvc := &MockRoleService{}
		h := NewUserHandler(userSvc, roleSvc)
		r := setupUserTestRouter()
		r.GET("/users/:email/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/users/fail@example.com/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}
