package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	sdk "github.com/fpt-corp/fptai-sdk-go"
)

func main() {
	username := "ftel"
	password := "ruandengming"
	application_code := "fa4e7de1a3e3251df910130f9cb1d375"

	client, err := sdk.NewClient(username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.SessionID())

	app := client.GetApp(application_code)
	ir, err := app.Recognize("Cảm ơn anh nhé")
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Train(); err != nil {
		log.Fatal("failed to train: ", err)
	}
	log.Info("done")
	fmt.Println(ir.Intent, ir.Confidence)

	intents, err := app.Intents()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(intents)

	intent, err := app.CreateIntent("goodbye", "Say Goodbye")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("new intent: ", intent)
}