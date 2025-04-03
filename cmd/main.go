package main

import (
	"certificate/internal/adapters"
	"certificate/internal/config"
	"certificate/internal/delivery"
	"certificate/internal/services"
	"log/slog"
	"os"
	"sync"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Ошибка загрузки конфигурации:", "error", err)
		os.Exit(1)
	}
	slog.Info("Успешная загрузка конфигурации")

	repo, err := adapters.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		slog.Error("Ошибка подключения к БД:", "error", err)
		os.Exit(1)
	}

	slog.Info("Успешное подключение к БД")

	api := adapters.NewPosterAPI(cfg.PosterToken)
	svc := services.NewRegistrationService(repo, api)

	bot, err := delivery.NewBot(cfg.BotToken, svc, cfg.BaseURL, cfg.Admins)
	if err != nil {
		slog.Error("Ошибка запуска бота:", "error", err)
		os.Exit(1)
	}

	server := delivery.NewHTTPServer(svc)
	server.ServeStaticFiles()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bot.Start()
	}()

	go func() {
		defer wg.Done()
		server.Start(":" + cfg.ServerPort)
	}()

	wg.Wait()
}
