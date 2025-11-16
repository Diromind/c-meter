package main

import (
	"log"
	"time"

	"backend/config"
	"backend/internal/bot"
	"backend/internal/bot/handlers"
	"backend/internal/database"

	tele "gopkg.in/telebot.v3"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Bot.Token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.RunMigrations("migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	pref := tele.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	handler := bot.NewBotHandler(db)
	menuHandler := handlers.NewMenuHandler(db)

	b.Handle("/start", handler.HandleStart)
	b.Handle("/help", handler.HandleHelp)
	b.Handle("/ping", handler.HandlePing)
	b.Handle("/get", handler.HandleGet)
	b.Handle("/today", handler.HandleGet)
	b.Handle("/record", handler.HandleRecord)
	b.Handle("/set_noon", handler.HandleSetNoon)
	b.Handle("/set_lang", handler.HandleSetLang)
	
	b.Handle("/menu", menuHandler.HandleMenu)
	b.Handle(&menuHandler.BtnLocations, menuHandler.HandleLocationsCallback)
	b.Handle(tele.OnCallback, menuHandler.HandleCallback)

	log.Println("Bot started successfully!")
	b.Start()
}
