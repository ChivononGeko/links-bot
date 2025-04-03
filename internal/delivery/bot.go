package delivery

import (
	"log"
	"log/slog"
	"strings"
	"time"

	"certificate/internal/ports"

	"github.com/tucnak/telebot"
)

type Bot struct {
	bot     *telebot.Bot
	svc     ports.RegistrationService
	baseURL string
	admins  map[int]struct{}
}

func NewBot(token string, svc ports.RegistrationService, baseURL string, adminIDs []int) (*Bot, error) {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	adminMap := make(map[int]struct{}, len(adminIDs))
	for _, id := range adminIDs {
		adminMap[id] = struct{}{}
	}

	return &Bot{bot: b, svc: svc, baseURL: baseURL, admins: adminMap}, nil
}

// Проверка, является ли пользователь админом
func (b *Bot) isAdmin(userID int) bool {
	_, exists := b.admins[userID]
	return exists
}

// Запуск бота
func (b *Bot) Start() {
	// Команда для генерации ссылки
	b.bot.Handle("/register", func(m *telebot.Message) {
		if !b.isAdmin(m.Sender.ID) {
			slog.Error("Попытка генерации токена, лицом без доступа", "ID", m.Sender.ID)
			b.bot.Send(m.Sender, "У вас нет прав для использования этой команды.")
			return
		}

		link, err := b.svc.GenerateUniqueLink(b.baseURL)
		if err != nil {
			slog.Error("Ошибка при создании ссылки", "error", err)
			b.bot.Send(m.Sender, "Ошибка при создании ссылки")
			return
		}
		b.bot.Send(m.Sender, "Ваша ссылка: "+link)
	})

	// Команда для проверки данных по токену
	b.bot.Handle("/check_token", func(m *telebot.Message) {
		if !b.isAdmin(m.Sender.ID) {
			slog.Error("Попытка првоерки токена, лицом без доступа", "ID", m.Sender.ID)
			b.bot.Send(m.Sender, "У вас нет прав для использования этой команды.")
			return
		}

		b.bot.Send(m.Sender, "Введите токен для проверки:")
	})

	// Обработчик сообщений (проверяем, ввел ли пользователь токен после команды)
	b.bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if !b.isAdmin(m.Sender.ID) {
			return
		}

		// Ищем данные по введенному токену
		usage, err := b.svc.GetTokenUsage(m.Text)
		if err != nil {
			slog.Error("Ошибка при поиске данных или токен не найден", "error", err)
			b.bot.Send(m.Sender, "Ошибка при поиске данных или токен не найден.")
			return
		}

		// Отправляем данные о пользователе
		response := "Данные по токену:\n"
		response += "👤 Имя: " + usage.Username + "\n"
		response += "📞 Телефон: " + usage.Phone

		b.bot.Send(m.Sender, response)
	})

	// Получение списка использованных токенов
	b.bot.Handle("/used_tokens", func(m *telebot.Message) {
		if !b.isAdmin(m.Sender.ID) {
			slog.Error("Попытка получении списка токенов, лицом без доступа", "ID", m.Sender.ID)
			b.bot.Send(m.Sender, "У вас нет прав для использования этой команды.")
			return
		}

		tokens, err := b.svc.GetUsedTokens()
		if err != nil {
			slog.Error("Ошибка при получении использованных токенов", "error", err)
			b.bot.Send(m.Sender, "Ошибка при получении использованных токенов.")
			return
		}

		if len(tokens) == 0 {
			b.bot.Send(m.Sender, "Нет использованных токенов.")
			return
		}

		response := "📌 *Список использованных токенов:*\n"
		for _, t := range tokens {
			escapedToken := strings.ReplaceAll(t.Token, "`", "\\`")
			response += "🔹 `" + escapedToken + "`\n"
		}

		b.bot.Send(m.Sender, response, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	})

	// Получение списка неиспользованных токенов
	b.bot.Handle("/unused_tokens", func(m *telebot.Message) {
		if !b.isAdmin(m.Sender.ID) {
			slog.Error("Попытка получении списка токенов, лицом без доступа", "ID", m.Sender.ID)
			b.bot.Send(m.Sender, "У вас нет прав для использования этой команды.")
			return
		}

		tokens, err := b.svc.GetUnusedTokens()
		if err != nil {
			slog.Error("Ошибка при получении неиспользованных токенов", "error", err)
			b.bot.Send(m.Sender, "Ошибка при получении неиспользованных токенов.")
			return
		}

		if len(tokens) == 0 {
			b.bot.Send(m.Sender, "Нет неиспользованных токенов.")
			return
		}

		response := "📌 *Список неиспользованных токенов:*\n"
		for _, t := range tokens {
			escapedToken := strings.ReplaceAll(t.Token, "`", "\\`")
			response += "🟢 `" + escapedToken + "`\n"
		}

		b.bot.Send(m.Sender, response, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	})

	log.Println("Бот запущен!")
	b.bot.Start()
}
