package main

import (
	"context"
	"fmt"
	"github.com/dimassfeb-09/openai_whatsapp/qrcode"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	gogpt "github.com/sashabaranov/go-gpt3"
	"go.mau.fi/whatsmeow"
	protoMsg "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

var client *whatsmeow.Client

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe {
			conversation := v.Message.GetConversation()
			if conversation != "" {

				rangeMsg := strings.Split(conversation, " ")
				for _, Msg := range rangeMsg {
					if Msg == ".openai" {

						openAI := gogpt.NewClient("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")
						ctx := context.Background()

						request := gogpt.CompletionRequest{
							Model:     "text-davinci-003",
							Prompt:    conversation,
							MaxTokens: 2000,
						}

						res, err := openAI.CreateCompletion(ctx, request)
						if err != nil {
							fmt.Println(err)
						}

						client.SendMessage(ctx, v.Info.Sender, "", &protoMsg.Message{
							Conversation: proto.String("Sedang mencari jawaban... [■■■■■■■■■□] 90%"),
						})

						response := strings.TrimSpace(res.Choices[0].Text)
						_, errMsg := client.SendMessage(ctx, v.Info.Sender, "", &protoMsg.Message{
							Conversation: proto.String(response),
						})
						if errMsg != nil {
							fmt.Println(errMsg.Error())
						}

					}

				}
			}
		}
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println("EXEC 1")
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		fmt.Println(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		// Create QR CODE
		qrcode.QrCode(client)
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
