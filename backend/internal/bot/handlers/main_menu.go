package handlers

import (
	"fmt"
	"log"
	"strings"

	"backend/internal/database"

	tele "gopkg.in/telebot.v3"
)

type MenuHandler struct {
	db           *database.DB
	BtnLocations tele.Btn
}

func NewMenuHandler(db *database.DB) *MenuHandler {
	menu := &tele.ReplyMarkup{}
	return &MenuHandler{
		db:           db,
		BtnLocations: menu.Data("📍 Locations", "menu_locations"),
	}
}

func (h *MenuHandler) HandleMenu(c tele.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(h.BtnLocations),
	)
	
	return c.Send("Main Menu:", menu)
}

func (h *MenuHandler) HandleCallback(c tele.Context) error {
	data := c.Callback().Data
	
	if len(data) > 0 && data[0] < 32 {
		data = data[1:]
	}
	
	log.Printf("Received callback data (cleaned): '%s'", data)
	
	if strings.HasPrefix(data, "nav:") {
		log.Printf("Matched nav: prefix")
		return h.HandleNavigationCallback(c)
	}
	
	if strings.HasPrefix(data, "add:") {
		log.Printf("Matched add: prefix")
		return c.Respond(&tele.CallbackResponse{Text: "Add item feature - coming soon!"})
	}
	
	log.Printf("Unknown callback action: '%s'", data)
	return c.Respond(&tele.CallbackResponse{Text: "Unknown action"})
}

func (h *MenuHandler) HandleLocationsCallback(c tele.Context) error {
	return h.showLocationsLevel(c, "")
}

func (h *MenuHandler) HandleNavigationCallback(c tele.Context) error {
	data := c.Callback().Data
	
	if len(data) > 0 && data[0] < 32 {
		data = data[1:]
	}
	
	if !strings.HasPrefix(data, "nav:") {
		return c.Respond(&tele.CallbackResponse{Text: "Invalid action"})
	}
	
	path := strings.TrimPrefix(data, "nav:")
	
	return h.showLocationsLevel(c, path)
}

func (h *MenuHandler) showLocationsLevel(c tele.Context, parentPath string) error {
	login := c.Sender().Username
	if login == "" {
		login = fmt.Sprintf("user_%d", c.Sender().ID)
	}
	
	items, err := h.db.GetUserCommonItemsAtLevel(login, parentPath)
	if err != nil {
		log.Printf("Error getting items at level: %v", err)
		return c.Respond(&tele.CallbackResponse{Text: "Error loading items"})
	}
	
	menu := &tele.ReplyMarkup{}
	var rows []tele.Row
	
	for _, item := range items {
		icon := "📁"
		if item.ProductUUID != nil {
			icon = "🍽️"
		}
		
		btnText := fmt.Sprintf("%s %s", icon, item.Name)
		btn := menu.Data(btnText, "nav:"+item.Path)
		rows = append(rows, menu.Row(btn))
	}
	
	btnAdd := menu.Data("➕ Add item", "add:"+parentPath)
	rows = append(rows, menu.Row(btnAdd))
	
	if parentPath != "" {
		parts := strings.Split(parentPath, ".")
		var parentPathStr string
		if len(parts) > 1 {
			parentPathStr = strings.Join(parts[:len(parts)-1], ".")
		}
		btnBack := menu.Data("⬅️ Back", "nav:"+parentPathStr)
		rows = append(rows, menu.Row(btnBack))
	}
	
	menu.Inline(rows...)
	
	title := "📍 Locations"
	if parentPath != "" {
		title = fmt.Sprintf("📂 %s", parentPath)
	}
	
	if len(items) == 0 {
		title += "\n\n(Empty - click ➕ to add items)"
	}
	
	err = c.Edit(title, menu)
	if err != nil {
		return c.Send(title, menu)
	}
	
	return c.Respond()
}

