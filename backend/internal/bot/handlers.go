package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"backend/internal/database"

	tele "gopkg.in/telebot.v3"
)

type BotHandler struct {
	db *database.DB
}

func NewBotHandler(db *database.DB) *BotHandler {
	return &BotHandler{db: db}
}

func (h *BotHandler) HandleStart(c tele.Context) error {
	return c.Send("Welcome to C-Meter! 👋\n\nUse /help to see available commands.")
}

func (h *BotHandler) HandleHelp(c tele.Context) error {
	helpText := `Available commands:

/start - Start the bot
/help - Show this help message
/ping - Check database connection and schema version
/get [days] - Get your entries (default: 1 day)
/today - Get today's entries
/record <name> <ccal> [proteins] [fats] [carbs] - Add a food record`

	return c.Send(helpText)
}

func (h *BotHandler) HandlePing(c tele.Context) error {
	version, err := h.db.GetLatestSchemaVersion()
	if err != nil {
		log.Printf("Error getting schema version: %v", err)
		return c.Send("❌ Database error: " + err.Error())
	}

	if version == "" {
		return c.Send("✅ Database connected\n📦 No migrations applied yet")
	}

	message := fmt.Sprintf("✅ Database connected\n📦 Schema version: %s", version)
	return c.Send(message)
}

func (h *BotHandler) HandleGet(c tele.Context) error {
	args := c.Args()
	days := 1
	
	if len(args) > 0 {
		var err error
		days, err = strconv.Atoi(args[0])
		if err != nil || days <= 0 {
			return c.Send("Please provide a valid number of days (positive integer)")
		}
	}

	login := c.Sender().Username
	if login == "" {
		login = fmt.Sprintf("user_%d", c.Sender().ID)
	}

	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(days) * 24 * time.Hour)

	records, err := h.db.GetRecordsByLoginAndTimeRange(login, startTime, endTime)
	if err != nil {
		log.Printf("Error getting records: %v", err)
		return c.Send("❌ Error fetching records: " + err.Error())
	}

	if len(records) == 0 {
		return c.Send("No records found for the last " + strconv.Itoa(days) + " days")
	}

	var result strings.Builder
	if days == 1 {
		result.WriteString("<b>Today's records:</b>\n\n")
	} else {
		result.WriteString(fmt.Sprintf("<b>Records for last %d days:</b>\n\n", days))
	}
	
	result.WriteString("<pre>")
	result.WriteString("date & time │      name       │ kcal\n")
	result.WriteString("────────────┼─────────────────┼─────\n")

	totalKcal := int64(0)
	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		
		product, err := h.db.GetProductByUUID(record.ProductUUID)
		if err != nil {
			log.Printf("Error getting product %s: %v", record.ProductUUID, err)
			continue
		}

		dateTime := record.CreatedAt.Format("02-01 15:04")
		name := product.Name
		if len(name) > 15 {
			name = name[:15]
		}
		
		padding := (15 - len(name)) / 2
		centeredName := fmt.Sprintf("%*s%s%*s", padding, "", name, 15-len(name)-padding, "")
		
		ccal := product.Ccal * record.Amount
		totalKcal += ccal

		line := fmt.Sprintf("%s │ %s │ %-4d\n", dateTime, centeredName, ccal)
		result.WriteString(line)
	}
	
	result.WriteString("</pre>")
	
	if days == 1 {
		result.WriteString(fmt.Sprintf("\n\n📋 <b>Total: %d kcal</b>", totalKcal))
	}

	return c.Send(result.String(), &tele.SendOptions{ParseMode: tele.ModeHTML})
}

func (h *BotHandler) HandleRecord(c tele.Context) error {
	args := c.Args()
	if len(args) < 2 {
		return c.Send("Usage: /record <name> <ccal> [proteins] [fats] [carbs]\nExample: /record \"Chicken Breast\" 165 31 3 0")
	}

	name := args[0]
	
	ccal, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || ccal <= 0 {
		return c.Send("Calories must be a positive number")
	}

	var proteins, fats, carbs int64

	if len(args) > 2 {
		proteins, err = strconv.ParseInt(args[2], 10, 64)
		if err != nil || proteins < 0 {
			return c.Send("Proteins must be a non-negative number")
		}
	}

	if len(args) > 3 {
		fats, err = strconv.ParseInt(args[3], 10, 64)
		if err != nil || fats < 0 {
			return c.Send("Fats must be a non-negative number")
		}
	}

	if len(args) > 4 {
		carbs, err = strconv.ParseInt(args[4], 10, 64)
		if err != nil || carbs < 0 {
			return c.Send("Carbs must be a non-negative number")
		}
	}

	product, err := h.db.InsertProduct(name, ccal, fats, proteins, carbs)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return c.Send("❌ Error creating product: " + err.Error())
	}

	login := c.Sender().Username
	if login == "" {
		login = fmt.Sprintf("user_%d", c.Sender().ID)
	}

	record, err := h.db.InsertRecord(product.UUID, 1, login)
	if err != nil {
		log.Printf("Error inserting record: %v", err)
		return c.Send("❌ Error creating record: " + err.Error())
	}

	message := fmt.Sprintf("✅ Recorded: %s\n📊 Calories: %d\nID: %s", name, ccal, record.UUID)
	return c.Send(message)
}

