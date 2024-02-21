package runcodeworker

import (
	"time"

	github.com/enteprise/etl-central/app/mainapp/database"
	modelmain github.com/enteprise/etl-central/app/mainapp/database/models"
	github.com/enteprise/etl-central/app/mainapp/messageq"
	"github.com/google/uuid"
)

/* Return errors to the logging console */
func WSLogError(logline string, msg modelmain.CodeRun) {

	uidstring := uuid.NewString()
	sendmsg := modelmain.LogsSend{
		CreatedAt:     time.Now().UTC(),
		UID:           uidstring,
		Log:           logline,
		LogType:       "error",
		EnvironmentID: msg.EnvironmentID,
		RunID:         msg.RunID,
	}

	messageq.MsgSend("coderunfilelogs."+msg.EnvironmentID+"."+msg.RunID, sendmsg)
	database.DBConn.Create(&sendmsg)

	sendmsg2 := modelmain.LogsSend{
		CreatedAt: time.Now().UTC(),
		UID:       uuid.NewString(),
		Log:       "Fail",
		LogType:   "action",
	}

	messageq.MsgSend("coderunfilelogs."+msg.EnvironmentID+"."+msg.RunID, sendmsg2)

}
