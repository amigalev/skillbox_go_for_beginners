package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type binanceResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

// Просто объявлена, но не проинициализирована
type wallet map[string]float64

// Проинициализирована
var db = map[int64]wallet{}

func main() {
	bot, err := tgbotapi.NewBotAPI("2115338446:AAHv0eB6cgv_onol2CipbKp9oozJXhpU3wo")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageArray := strings.Split(update.Message.Text, " ")

		switch messageArray[0] {
		case "ADD":
			summ, err := strconv.ParseFloat(messageArray[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка конвертации значения "+messageArray[2]))
				continue
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = wallet{}
			}
			db[update.Message.Chat.ID][messageArray[1]] += summ

			result := fmt.Sprintf("Баланс: %s %f", messageArray[1], db[update.Message.Chat.ID][messageArray[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, result))
		case "SUB":
			summ, err := strconv.ParseFloat(messageArray[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка конвертации значения "+messageArray[2]))
				continue
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = wallet{}
			}
			db[update.Message.Chat.ID][messageArray[1]] -= summ

			result := fmt.Sprintf("Баланс: %s %f", messageArray[1], db[update.Message.Chat.ID][messageArray[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, result))
		case "DEL":
			delete(db[update.Message.Chat.ID], messageArray[1])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта удалена"))
		case "SHOW":
			msg := "Баланс:\n"
			var usdSumm float64
			for key, value := range db[update.Message.Chat.ID] {
				coinPrice, err := getPrice(key)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка конвертации валюты: "+err.Error()))
				}

				usdSumm += value * coinPrice
				msg += fmt.Sprintf("%s: %f [%.2f]\n", key, value, value*coinPrice)
			}
			msg += fmt.Sprintf("Сумма: %.2f", usdSumm)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная комманда"))
		}
		// msg.ReplyToMessageID = update.Message.MessageID

	}
}

// api.binance.com/api/v3/ticker/price?symbol=BTCUSDT

// XRP
// ETH

func getPrice(coin string) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf("http://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var jsonResponce binanceResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResponce)
	if err != nil {
		return
	}

	if jsonResponce.Code != 0 {
		err = errors.New("Некорректная валюта")
	}

	price = jsonResponce.Price
	return
}

// Не уйти в минус
// Нельзя вводить неверную валюту
// Выводить надо в рублях
