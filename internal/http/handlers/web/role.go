package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/user"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	svc user.RoleService
}

func NewRoleHandler(svc user.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	roles, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/roles/index.html", gin.H{
		"title":           "Roles",
		"roles":           roles,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *RoleHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/roles/form.html", gin.H{
		"title":  "Add Role",
		"action": "/roles",
	})
}

func (h *RoleHandler) Create(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")

	_, err := h.svc.CreateRole(id, name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/roles")
}

func (h *RoleHandler) View(c *gin.Context) {
	id := c.Param("id")

	role, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		c.String(http.StatusNotFound, "Role not found")
		return
	}

	c.HTML(http.StatusOK, "pages/roles/view.html", gin.H{
		"title": "View Role",
		"role":  role,
	})
}

func (h *RoleHandler) EditForm(c *gin.Context) {
	id := c.Param("id")

	role, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		c.String(http.StatusNotFound, "Role not found")
		return
	}

	c.HTML(http.StatusOK, "pages/roles/form.html", gin.H{
		"title":  "Edit Role",
		"action": "/roles/" + id,
		"role":   role,
	})
}

func (h *RoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")

	_, err := h.svc.UpdateRole(id, name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/roles")
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.DeleteRole(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *RoleHandler) DeleteConfirm(c *gin.Context) {
	id := c.Param("id")

	role, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		c.String(http.StatusNotFound, "Role not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      role.Name,
		"DeleteURL": "/roles/" + id,
		"RowID":     "role-row-" + id,
	})
}
