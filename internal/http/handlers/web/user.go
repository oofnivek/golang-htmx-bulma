package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/service"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
	roleSvc service.RoleService
}

func NewUserHandler(userSvc service.UserService, roleSvc service.RoleService) *UserHandler {
	return &UserHandler{userSvc: userSvc, roleSvc: roleSvc}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "email")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	search := c.Query("search")
	tz := c.DefaultQuery("tz", "UTC")

	users, total, err := h.userSvc.ListPaged(page, pageSize, sortBy, sortOrder, search)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/users/index.html", gin.H{
		"title":           "Users",
		"users":           users,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
		"search":          search,
		"tz":              tz,
	})
}

func (h *UserHandler) CreateForm(c *gin.Context) {
	roles, err := h.roleSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/users/form.html", gin.H{
		"title":  "Add User",
		"action": "/users",
		"roles":  roles,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	isEnabled := c.PostForm("is_enabled") == "on"
	user := &model.User{
		Email:       c.PostForm("email"),
		FirstName:   c.PostForm("first_name"),
		LastName:    c.PostForm("last_name"),
		Mobile:      c.PostForm("mobile"),
		Designation: c.PostForm("designation"),
		Department:  c.PostForm("department"),
		IsEnabled:   isEnabled,
		RoleID:      c.PostForm("role_id"),
		Password:    c.PostForm("password"),
		ConfirmPass: c.PostForm("confirm_password"),
	}

	err := h.userSvc.CreateUser(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")
}

func (h *UserHandler) EditForm(c *gin.Context) {
	email := c.Param("email")
	user, err := h.userSvc.GetByEmail(email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	roles, err := h.roleSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/users/form.html", gin.H{
		"title":  "Edit User",
		"action": "/users/" + email,
		"user":   user,
		"roles":  roles,
	})
}

func (h *UserHandler) PasswordFields(c *gin.Context) {
	c.HTML(http.StatusOK, "partials/users/password_fields.html", nil)
}

func (h *UserHandler) View(c *gin.Context) {
	email := c.Param("email")
	tz := c.DefaultQuery("tz", "UTC")
	user, err := h.userSvc.GetByEmail(email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.HTML(http.StatusOK, "pages/users/view.html", gin.H{
		"title": "View User",
		"user":  user,
		"tz":    tz,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	email := c.Param("email")
	isEnabled := c.PostForm("is_enabled") == "on"
	
	user := &model.User{
		Email:       email,
		FirstName:   c.PostForm("first_name"),
		LastName:    c.PostForm("last_name"),
		Mobile:      c.PostForm("mobile"),
		Designation: c.PostForm("designation"),
		Department:  c.PostForm("department"),
		IsEnabled:   isEnabled,
		RoleID:      c.PostForm("role_id"),
		Password:    c.PostForm("password"),
		ConfirmPass: c.PostForm("confirm_password"),
	}

	err := h.userSvc.UpdateUser(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")
}

func (h *UserHandler) Delete(c *gin.Context) {
	email := c.Param("email")
	err := h.userSvc.DeleteUser(email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *UserHandler) DeleteConfirm(c *gin.Context) {
	email := c.Param("email")
	user, err := h.userSvc.GetByEmail(email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      user.FirstName + " " + user.LastName + " (" + user.Email + ")",
		"DeleteURL": "/users/" + email,
		"RowID":     "user-row-" + view.SafeID(email),
	})
}
