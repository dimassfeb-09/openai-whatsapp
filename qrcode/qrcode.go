package qrcode

import (
	"context"
	"fmt"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"log"
)

func QrCode(client *whatsmeow.Client) {
	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		fmt.Println("Error Connection")
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr.png")
			if err != nil {
				log.Println(err)
			}
		} else {
			fmt.Println("Login event:", evt.Event)
		}

	}
}
