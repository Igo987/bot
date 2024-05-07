package bot

import (
	"context"
	"fmt"
	"github/Igo87/crypt/config"
	"github/Igo87/crypt/models"
	"github/Igo87/crypt/pkg/logger"
	"strconv"

	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Bot *tgbotapi.BotAPI
}

type Command struct {
	ChatID  int64
	Command string
}

func NewBotAPI(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		err = fmt.Errorf("failed to create bot: %w", err)

		return nil, err
	}

	return &Bot{Bot: bot}, nil
}

// Start initializes and starts the connection with Telegram API.
func Start() (*Bot, error) {
	bot, err := NewBotAPI(config.Cfg.GetToken())
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return bot, nil
}

// SendMessage sends data to a bot.
func (b *Bot) SendData(ctx context.Context, data models.Currencies, l models.Currencies) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.Bot.GetUpdatesChan(u)

	var someMap = make(map[int]string)

	interval := 20

	for {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				logger.Log.Error("failed to get updates")

				return
			}

			cmd := update.Message.Command()
			if cmd != "" {
				if update.Message.CommandArguments() != "" {
					minutes, err := strconv.ParseFloat(update.Message.CommandArguments(), 64)
					if err != nil {
						_, err := b.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, `Неверный формат команды
 введите число (в минутах)`))
						if err != nil {
							logger.Log.Error("failed to send message: %v", err)
						}

						continue
					}

					interval = int(minutes)
				}

				someMap[int(update.Message.Chat.ID)] = cmd
			}

		case <-ticker.C:
			for chatID, cmd := range someMap {
				if err := b.HandleCommand(cmd, int64(chatID), data, l, ctx); err != nil {
					delete(someMap, chatID)
				}
			}
		}
	}
}

func (b *Bot) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	tgMessage, err := b.Bot.Send(msg)
	if err != nil {
		return tgMessage, fmt.Errorf("failed to send message: %w", err)
	}

	return tgMessage, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.Send(msg)
	return err
}

func (b *Bot) SendMessageBTC(chatID int64, data models.Crypto) error {
	btc := data.Data.Bitcoin
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		`\f Валюта: %s,\n текущая цена: %.10f,\n изменение за 1 час: %.10f,
		\n изменение за 24 ч.: %.8f,\n изменение за 7 дн.: %.8f,\n изменение за 30 дн.: %.8f`,
		btc.Name, btc.Quote.Rub.Price,
		btc.Quote.Rub.PercentChange1H, btc.Quote.Rub.PercentChange24H,
		btc.Quote.Rub.PercentChange7D, btc.Quote.Rub.PercentChange30D,
	))
	_, err := b.Send(msg)
	return err
}

func (b *Bot) SendMessageBTCByYesterday(chatID int64, data models.Currencies) error {
	btc := data[0]
	yesterday := time.Now().AddDate(0, 0, -1).Format("01/02/2006")
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"\f Значения за %s :\n Валюта: %s,\n Процент изменения: %.10f,\n Мин.значение: %.10f,\n Макс.значение: %.10f",
		yesterday, btc.Name, btc.Percent, btc.Min, btc.Max,
	))
	_, err := b.Send(msg)
	return err
}

func (b *Bot) SendMessageETHByYesterday(chatID int64, data models.Currencies) error {
	eth := data[1]
	yesterday := time.Now().AddDate(0, 0, -1).Format("01/02/2006")
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"  Значения за %s :\n Валюта: %s,\n Процент изменения: %.10f,\n Мин.значение: %.10f,\n Макс.значение: %.10f",
		yesterday, eth.Name, eth.Percent, eth.Min, eth.Max,
	))
	_, err := b.Send(msg)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) SendMessageETH(chatID int64, data models.Crypto) error {
	eth := data.Data.Ethereum

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		`Валюта: %s,\n текущая цена: %.10f,\n изменение за 1 час: %.10f,
		\n изменение за 24 ч.: %.8f,\n изменение за 7 дн.: %.8f,\n изменение за 30 дн.: %.8f`,
		eth.Name, eth.Quote.Rub.Price,
		eth.Quote.Rub.PercentChange1H, eth.Quote.Rub.PercentChange24H,
		eth.Quote.Rub.PercentChange7D, eth.Quote.Rub.PercentChange30D,
	))
	_, err := b.Send(msg)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) BTHExtr(chatID int64, data models.Currencies) error {
	btc := data[0]
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"  Валюта: %s,\n Изменение в процентах за час: %.10f,\n Мин.значение: %.10f,\n Макс.значение: %.10f",
		btc.Name, btc.Percent, btc.Min, btc.Max,
	))
	_, err := b.Send(msg)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) ETHExtr(chatID int64, data models.Currencies) error {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"  Валюта: %s,\n Изменение в процентах за час: %.10f,\n Мин.значение: %.10f,\n Макс.значение: %.10f",
		data[1].Name, data[1].Percent, data[1].Min, data[1].Max,
	))

	_, err := b.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// HandleCommand processes commands.
func (b *Bot) HandleCommand(update string, chaID int64, data models.Currencies, last models.Currencies, ctx context.Context) error {

	var msg string

	var err error

	switch update {
	case "start":
		msg = `Привет, я бот для криптовалют. Напиши /btc или /eth`
	case "btc":
		err = b.BTHExtr(chaID, data)
	case "eth":
		err = b.ETHExtr(chaID, data)
	case "eth_y":
		err = b.SendMessageETHByYesterday(chaID, last)
	case "btc_y":
		err = b.SendMessageBTCByYesterday(chaID, last)
	case "stop":
		err := b.SendMessage(chaID, "Бот остановлен")
		if err != nil {
			return err
		}

		b.Bot.StopReceivingUpdates()

		return nil

	default:
		msg = "Я пока не знаю такую команду"
	}

	if err != nil {
		return err
	}

	if msg != "" {
		return b.SendMessage(chaID, msg)
	}

	return nil
}
