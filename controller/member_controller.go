package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vickyaruldoss/ims/model"
	"github.com/vickyaruldoss/ims/repository"
	"github.com/vickyaruldoss/ims/service"
)

type MemberController struct {
	service service.MemberService
}

func NewMemberController(svc service.MemberService) *MemberController {
	return &MemberController{service: svc}
}

// CreateMember godoc
// @Summary      Create a member
// @Description  Add a new member to the institution
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        member  body      model.CreateMemberRequest  true  "Member payload"
// @Success      201     {object}  model.Member
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /members [post]
func (c *MemberController) CreateMember(ctx *gin.Context) {
	var req model.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	member, err := c.service.CreateMember(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, member)
}

// GetMember godoc
// @Summary      Get a member by ID
// @Description  Retrieve a single member by their ID
// @Tags         members
// @Produce      json
// @Param        id   path      string  true  "Member ID"
// @Success      200  {object}  model.Member
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /members/{id} [get]
func (c *MemberController) GetMember(ctx *gin.Context) {
	id := ctx.Param("id")
	member, err := c.service.GetMember(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, member)
}

// GetAllMembers godoc
// @Summary      List all members
// @Description  Retrieve all members in the institution
// @Tags         members
// @Produce      json
// @Success      200  {array}   model.Member
// @Failure      500  {object}  map[string]string
// @Router       /members [get]
func (c *MemberController) GetAllMembers(ctx *gin.Context) {
	members, err := c.service.GetAllMembers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, members)
}

// UpdateMember godoc
// @Summary      Update a member
// @Description  Update one or more fields of an existing member
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        id      path      string                     true  "Member ID"
// @Param        member  body      model.UpdateMemberRequest  true  "Fields to update"
// @Success      200     {object}  model.Member
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /members/{id} [put]
func (c *MemberController) UpdateMember(ctx *gin.Context) {
	id := ctx.Param("id")
	var req model.UpdateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	member, err := c.service.UpdateMember(id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, member)
}

// DeleteMember godoc
// @Summary      Delete a member
// @Description  Remove a member from the institution by ID
// @Tags         members
// @Param        id   path  string  true  "Member ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /members/{id} [delete]
func (c *MemberController) DeleteMember(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteMember(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
