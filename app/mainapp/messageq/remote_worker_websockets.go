package messageq

import "github.com/enteprise/etl-central/app/mainapp/database/models"

var WebsocketRWChannel = make(chan models.WSChannelMessage)
