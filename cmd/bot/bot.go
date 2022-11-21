package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"upgrade/internal/models"

	"gopkg.in/telebot.v3"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
}

func CreateBot() UpgradeBot {

	db, err := gorm.Open(mysql.Open("root:yes@tcp(localhost:3306)/go?parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Data base connecting error %v", err)
	}

	token := "5777491063:AAFlxsonvdC4YkPicvlx5n0NPlztVJ7mGdY"

	upgradeBot := UpgradeBot{
		Bot:   InitBot(token),
		Users: &models.UserModel{Db: db},
	}

	return upgradeBot
}

func hello(w http.ResponseWriter, r *http.Request) {

	body_byte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Body reading error: %v", err)
	}

	body_string := string(body_byte)

	body_trimmed := strings.TrimPrefix(body_string, "text=")

	body_decoded, err := url.QueryUnescape(body_trimmed)
	if err != nil {
		log.Fatalf("Body decoding error: %v", err)
	}

	fmt.Println(body_decoded)

	upgradeBot := CreateBot()
	upgradeBot.SendAll(body_decoded)
}

func send(chat_id int, text string) {

	chat_str := strconv.Itoa(chat_id)

	postBody, _ := json.Marshal(map[string]string{
		"chat_id": chat_str,
		"text":    text,
	})
	responseBody := bytes.NewBuffer(postBody)

	upgradeBot := CreateBot()
	resp, err := http.Post("https://api.telegram.org/bot"+upgradeBot.Bot.Token+"/sendMessage", "application/json", responseBody)

	if err != nil {
		log.Fatalf("Request sending error to Telegram: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Body reading error: %v", err)
	}
	sb := string(body)
	log.Println(sb)
}

func (bot *UpgradeBot) SendAll(text string) {

	usersId, err := bot.Users.FindAll()

	if err != nil {
		log.Printf("User get error %v", err)
	}

	for i := 0; i < len(usersId); i++ {
		send(usersId[i], text)
	}
}

func (bot *UpgradeBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("User get error %v", err)
	}

	if existUser == nil {
		err := bot.Users.Create(newUser)

		if err != nil {
			log.Printf("User creating error %v", err)
		}
	}

	return ctx.Send("Hello, " + ctx.Sender().FirstName)
}

func InitBot(token string) *telebot.Bot {

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Bot initialization error %v", err)
	}

	return bot
}

func Listen() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8090", nil)
}
