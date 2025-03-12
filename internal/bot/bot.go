package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/database"
	"github.com/4Amangel1/smart-house-automate/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api             *tgbotapi.BotAPI
	repo            *database.Repository
	authorizedUsers map[int64]bool
	logger          *log.Logger
	stopChan        chan struct{}
	config          config.TelegramBotConfig
}

func NewBot(cfg config.TelegramBotConfig, repo *database.Repository, logger *log.Logger) (*Bot, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	logger.Printf("Bot starting with token length: %d", len(cfg.Token))
	logger.Printf("Authorized users: %v", cfg.AuthorizedUserIDs)

	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot := &Bot{
		api:             api,
		repo:            repo,
		authorizedUsers: make(map[int64]bool),
		logger:          logger,
		stopChan:        make(chan struct{}),
		config:          cfg,
	}

	for _, id := range cfg.AuthorizedUserIDs {
		bot.authorizedUsers[id] = true
		logger.Printf("Added authorized user from config: %d", id)
	}

	return bot, nil
}

func (b *Bot) Start() error {

	b.logger.Printf("Bot starting with token length: %d", len(b.config.Token))
	b.logger.Printf("Authorized users map: %v", b.authorizedUsers)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	b.logger.Printf("Authorized users: %d", len(b.authorizedUsers))
	b.logger.Printf("Bot %s started successfully", b.api.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		b.logger.Printf("Received message from user %d in chat %d: %s",
			userID, chatID, update.Message.Text)

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
		}
	}

	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := message.From.ID

	b.logger.Printf("Received command from user %d: %s", userID, message.Command())

	switch message.Command() {
	case "start":
		if len(b.authorizedUsers) == 0 {
			b.authorizedUsers[userID] = true
			b.logger.Printf("Added first user as authorized: %d", userID)
		}

		text := "Добро пожаловать в систему управления умным домом!\n\n" +
			"Доступные команды:\n" +
			"/status - текущие показания датчиков\n" +
			"/history [sensor_id] - история показаний конкретного датчика"

		b.sendMessage(chatID, text)

	case "status":
		if !b.isAuthorized(userID) {
			b.sendMessage(chatID, "У вас нет прав для использования этой команды.")
			return
		}

		readings, err := b.repo.GetLatestReadings()
		if err != nil {
			b.logger.Printf("Error getting latest readings: %v", err)
			b.sendMessage(chatID, "Ошибка при получении данных.")
			return
		}

		if len(readings) == 0 {
			b.sendMessage(chatID, "Нет доступных показаний.")
			return
		}

		var msgBuilder strings.Builder
		msgBuilder.WriteString("📊 *Текущие показания датчиков:*\n\n")

		for _, reading := range readings {
			formattedValue := b.formatSensorValue(reading)
			msgBuilder.WriteString(fmt.Sprintf("🔹 *%s* (%s):\n%s\n\n",
				reading.SensorID,
				b.translateSensorType(reading.SensorType),
				formattedValue,
			))
		}

		b.sendMarkdownMessage(chatID, msgBuilder.String())

	case "history":
		if !b.isAuthorized(userID) {
			b.sendMessage(chatID, "У вас нет прав для использования этой команды.")
			return
		}

		args := message.CommandArguments()
		if args == "" {
			b.sendMessage(chatID, "Пожалуйста, укажите ID датчика.\nПример: /history temp_sensor_1")
			return
		}

		sensorID := args
		readings, err := b.repo.GetReadingHistory(sensorID, 10)
		if err != nil {
			b.logger.Printf("Error getting reading history: %v", err)
			b.sendMessage(chatID, "Ошибка при получении истории показаний.")
			return
		}

		if len(readings) == 0 {
			b.sendMessage(chatID, fmt.Sprintf("История показаний для датчика %s не найдена.", sensorID))
			return
		}

		var msgBuilder strings.Builder
		msgBuilder.WriteString(fmt.Sprintf("📈 *История показаний датчика %s:*\n\n", sensorID))

		for i, reading := range readings {
			formattedTime := reading.Timestamp.Format("02.01.2006 15:04:05")
			formattedValue := b.formatSensorValue(reading)

			msgBuilder.WriteString(fmt.Sprintf("%d. *%s*\n%s\n\n",
				i+1,
				formattedTime,
				formattedValue,
			))
		}

		b.sendMarkdownMessage(chatID, msgBuilder.String())

	default:
		b.sendMessage(chatID, "Неизвестная команда. Используйте /start для получения списка доступных команд.")
	}
}

func (b *Bot) LoadAuthorizedUsers(userIDsString string) {
	userIDs := strings.Split(userIDsString, ",")

	for _, idStr := range userIDs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			b.logger.Printf("Invalid user ID: %s", idStr)
			continue
		}

		b.authorizedUsers[id] = true
		b.logger.Printf("Added authorized user: %d", id)
	}
}

func (b *Bot) StartNotificationService() {
	ticker := time.NewTicker(b.config.AlertCheckInterval)
	defer ticker.Stop()

	b.logger.Printf("Starting notification service with interval %s", b.config.AlertCheckInterval)

	for {
		select {
		case <-ticker.C:
			b.checkTemperatureAlerts()
			b.checkMotionAlerts()
			b.checkAirQualityAlerts()
		case <-b.stopChan:
			b.logger.Println("Stopping notification service")
			return
		}
	}
}

func (b *Bot) checkTemperatureAlerts() {
	readings, err := b.repo.GetLatestReadingsByType("temperature")
	if err != nil {
		b.logger.Printf("Error getting temperature readings: %v", err)
		return
	}

	for _, reading := range readings {
		if value, ok := reading.Value.Data.(float64); ok {
			if value > b.config.TemperatureAlertThreshold {
				message := fmt.Sprintf("🔥 *ВНИМАНИЕ!* Высокая температура (%0.1f°C) на датчике %s",
					value, reading.SensorID)
				b.notifyAllUsers(message)
			}
		}
	}
}

func (b *Bot) checkMotionAlerts() {
	readings, err := b.repo.GetLatestReadingsByType("motion")
	if err != nil {
		b.logger.Printf("Error getting motion readings: %v", err)
		return
	}

	for _, reading := range readings {
		if value, ok := reading.Value.Data.(bool); ok && value {
			message := fmt.Sprintf("👤 *Обнаружено движение* на датчике %s", reading.SensorID)
			b.notifyAllUsers(message)
		}
	}
}

func (b *Bot) checkAirQualityAlerts() {
	readings, err := b.repo.GetLatestReadingsByType("air_quality")
	if err != nil {
		b.logger.Printf("Error getting air quality readings: %v", err)
		return
	}

	for _, reading := range readings {
		if valueMap, ok := reading.Value.Data.(map[string]interface{}); ok {
			if co2, co2ok := valueMap["co2"].(float64); co2ok && co2 > 1000 {
				message := fmt.Sprintf("⚠️ *ВНИМАНИЕ!* Высокий уровень CO₂ (%0.0f ppm) на датчике %s",
					co2, reading.SensorID)
				b.notifyAllUsers(message)
			}
		}
	}
}

func (b *Bot) notifyAllUsers(message string) {
	for userID := range b.authorizedUsers {
		b.sendMarkdownMessage(userID, message)
	}
}

func (b *Bot) isAuthorized(userID int64) bool {
	return b.authorizedUsers[userID]
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Printf("Error sending message to %d: %v", chatID, err)
	}
}

func (b *Bot) sendMarkdownMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Printf("Error sending markdown message to %d: %v", chatID, err)
		b.sendMessage(chatID, text)
	}
}

func (b *Bot) formatSensorValue(reading models.SensorData) string {
	switch reading.SensorType {
	case "temperature":
		value, ok := reading.Value.Data.(float64)
		if ok {
			return fmt.Sprintf("🌡 Температура: *%.1f°C*", value)
		}
	case "motion":
		value, ok := reading.Value.Data.(bool)
		if ok {
			if value {
				return "🔴 *Обнаружено движение*"
			} else {
				return "🟢 *Движение не обнаружено*"
			}
		}
	case "air_quality":
		if valueMap, ok := reading.Value.Data.(map[string]interface{}); ok {
			co2, co2ok := valueMap["co2"].(float64)
			nh3, nh3ok := valueMap["nh3"].(float64)

			if co2ok && nh3ok {
				return fmt.Sprintf("💨 CO₂: *%.0f ppm*\n💨 NH₃: *%.0f ppm*", co2, nh3)
			}
		}
	}

	return fmt.Sprintf("%v", reading.Value.Data)
}

func (b *Bot) translateSensorType(sensorType string) string {
	switch sensorType {
	case "temperature":
		return "Температура"
	case "motion":
		return "Движение"
	case "air_quality":
		return "Качество воздуха"
	default:
		return sensorType
	}
}

func (b *Bot) Stop() {
	b.logger.Println("Stopping bot...")
	close(b.stopChan)
}
