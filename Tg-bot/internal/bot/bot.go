package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"sync"
	"tg-bot/internal/fileManager"
	"tg-bot/internal/network"
)

type Bot interface {
	Start()
	SendMessage(chatId int64, message string)
	Stop()
}
type TgBot struct {
	*tgbotapi.BotAPI
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewTgBot(filePath string) *TgBot {

	tgBot := &TgBot{}

	config, err := fileManager.LoadConfig(filePath)
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal("Telegram token is empty")
	}

	botCommands, err := fileManager.CommandConverter(config.Commands)
	if err != nil {
		log.Fatal("Failed to convert commands: ", err)
	}

	setCommands := tgbotapi.NewSetMyCommands(botCommands...)
	_, err = bot.Request(setCommands)
	if err != nil {
		log.Fatal("Failed to set commands: ", err)
	}

	if config.ServerURL == "" {
		log.Fatal("Server URL is empty")
	}

	network.ServerURL = config.ServerURL

	log.Println("The config was uploaded successfully")

	tgBot.BotAPI = bot

	return tgBot
}

func (bot *TgBot) Start() {
	log.Println("The bot has been successfully launched...")

	bot.stopChan = make(chan struct{})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go MainController(updates, bot, bot.stopChan)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	bot.Stop()
}

func (bot *TgBot) SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (bot *TgBot) Stop() {
	log.Println("The bot is shutting down its work...")
	bot.StopReceivingUpdates()
	close(bot.stopChan)
	bot.wg.Wait()
	log.Println("The bot was successfully stopped!")
}

func (bot *TgBot) GetBotCommand() ([]tgbotapi.BotCommand, error) {
	commands, err := bot.GetMyCommands()
	if err != nil {
		return nil, err
	}
	return commands, nil
}
