package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovery()
	r.Run(":3434")
}
func Recovery() {
	ReadUser()
	botResult, err := tgbotapi.NewBotAPI("7016002303:AAHVdt5mao978Da8UtOBe_J91nL9SlRQ9cQ")

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	updates := tgbotapi.NewUpdate(0)
	update, _ := botResult.GetUpdatesChan(updates)

	for update := range update {
		if update.Message.IsCommand() {
			if update.Message.Command() == "reset" {
				for _, item := range SignUpSlice {
					if item.TelegramLogin == update.Message.Chat.UserName {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter New Password")
						botResult.Send(msg)
					}

				}
			}
		} else {

			for index, item := range SignUpSlice {

				if item.TelegramLogin == update.Message.Chat.UserName {

					SignUpSlice[index].Password = update.Message.Text

				}

			}

			WriteUser()

		}
	}
}
func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		WriteUser()
	}
}

func WriteUser() {
	marsheledData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", marsheledData, 0664)

}
func ReadUser() {
	readedByte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readedByte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
