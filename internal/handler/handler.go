package handler

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/user/fintech-expense-tracker/internal/service"
)

type Handler struct {
	svc service.Service
}

func NewHandler(s service.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) CreateUser(c *gin.Context) {
	var input struct { Name string `json:"name" binding:"required"`; Email string `json:"email" binding:"required,email"` }
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid user input: "+err.Error())
		return
	}
	u, err := h.svc.CreateUser(input.Name, input.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) CreateGroup(c *gin.Context) {
	var input struct { Name string `json:"name" binding:"required"` }
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid group input: "+err.Error())
		return
	}
	g, err := h.svc.CreateGroup(input.Name)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create group")
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) AddUserToGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Group ID must be a number")
		return
	}
	var input struct { UserID uint `json:"user_id" binding:"required"` }
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid member input: "+err.Error())
		return
	}
	if err := h.svc.AddUserToGroup(uint(id), input.UserID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add user to group")
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) AddExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Group ID must be a number")
		return
	}
	var input struct {
		PaidByID       uint           `json:"paid_by_id" binding:"required"`
		Description    string         `json:"description" binding:"required"`
		Amount         int64          `json:"amount" binding:"required"`
		IdempotencyKey string         `json:"idempotency_key" binding:"required"`
		Splits         map[uint]int64 `json:"splits" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid expense input: "+err.Error())
		return
	}
	e, err := h.svc.AddExpense(uint(id), input.PaidByID, input.Description, input.Amount, input.Splits, input.IdempotencyKey)
	if err != nil {
		code := "EXPENSE_ERROR"
		msg := err.Error()
		if parts := strings.Split(msg, ": "); len(parts) > 1 {
			code = parts[0]
			msg = parts[1]
		}
		RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *Handler) GetBalances(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Group ID must be a number")
		return
	}
	b, err := h.svc.GetBalances(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "EMPTY_GROUP") {
			RespondWithError(c, http.StatusBadRequest, "EMPTY_GROUP", "Group has no members")
			return
		}
		RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch balances")
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) GetSettlements(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Group ID must be a number")
		return
	}
	s, err := h.svc.GetSettlements(uint(id))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to compute settlements")
		return
	}
	c.JSON(http.StatusOK, s)
}
