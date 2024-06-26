package handlers

// The handler layer is responsible for parsing a request,
// calling out the relevant service and then returning a response to the caller.

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nazerkhh/shop/models"
	"github.com/nazerkhh/shop/services"
)

// ItemHandler is an interface for itemHandler struct
type ItemHandler interface {
	ListItems(c *gin.Context)
	FindItem(c *gin.Context)
	CreateItem(c *gin.Context)
	UpdateItem(c *gin.Context)
	DeleteItem(c *gin.Context)
}

type itemHandler struct {
	s *services.ItemService
}

// NewItemHandler is a constructor for ItemHandler
func NewItemHandler(s services.ItemService) ItemHandler {
	if s == nil {
		log.Fatal("Failed to initialize item handler, service is nil")
		return nil
	}
	var p = itemHandler{s: &s}
	return &p
}

func (h *itemHandler) SetItemService(s services.ItemService) {
	h.s = &s
}

func (h *itemHandler) GetItemService() services.ItemService {
	if h.s == nil {
		log.Fatal("Failed to get item service, it is nil")
		return nil
	}

	return *h.s
}

// GET /items - List all items
func (h *itemHandler) ListItems(c *gin.Context) {
	s := h.GetItemService()
	items, err := s.ListItems()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /items/:id - Find an item
func (h *itemHandler) FindItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bad id value: %s", err)})
		return
	}
	s := h.GetItemService()
	item, found, err := s.FindItem(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item not found"})
		return
	}
	if found {
		c.JSON(http.StatusOK, gin.H{"data": item})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"data": nil})
}

// POST /items - Create a new item
func (h *itemHandler) CreateItem(c *gin.Context) {
	// Validate input
	var input models.CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create an item
	s := h.GetItemService()
	item, err := s.CreateItem(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// PATCH /items/:id - Update an item
func (h *itemHandler) UpdateItem(c *gin.Context) {
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bad id value: %s", e)})
		return
	}
	// Validate input
	var input models.UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := h.GetItemService()
	item, err := s.UpdateItem(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DELETE /items/:id - Delete an item
func (h *itemHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bad id value: %s", err)})
		return
	}
	s := h.GetItemService()
	if err := s.DeleteItem(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
