package matrix

import (
	"bitsnthings.dev/overlord/src/log"
	"bitsnthings.dev/overlord/src/state/conf"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func parseMatrixMessage(src mautrix.EventSource, evt *event.Event) {
	log.PrintLog(log.INFO, "<%s> %s", evt.Sender, evt.Content.AsMessage().Body)
}

func Setup(conf conf.Config) {
	client, err := mautrix.NewClient(conf.MatrixCreds.Homeserver, "", "")
	if err != nil {
		log.PrintLog(log.ERROR, "Could not connect to matrix server! %s", err)
		conf.EnableMatrix = false
		return
	}
	_, err = client.Login(&mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: mautrix.IdentifierTypeUser, User: conf.MatrixCreds.Username},
		Password:         conf.MatrixCreds.Password,
		StoreCredentials: false,
	})
	if err != nil {
		log.PrintLog(log.ERROR, "Could not login to matrix server! %s", err)
		return
	}
	conf.MatrixCleint = *client
	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, parseMatrixMessage)
}
