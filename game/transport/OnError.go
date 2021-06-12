package transport

import (
	"github.com/nentenpizza/werewolves/wserver"
	log "github.com/sirupsen/logrus"
)

func (g *game) OnError(err error, ctx *wserver.Context) {
	client, ok := ctx.Get("client").(*Client)
	if ok {
		if err != nil {
			if client.Player != nil {
				Logger.WithFields(log.Fields{
					"client_name": client.Name,
					"client_role": client.Role,
				}).Error(err.Error())
			} else {
				Logger.WithFields(log.Fields{
					"client": client,
				}).Error(err.Error())
			}

			e, k := err.(*ServerError)
			if k {
				client.WriteJSON(EventErr{Type: ctx.EventType(), Data: ctx.Data(), Error: e})
			}
		}
		return
	}
	if err != nil {
		Logger.WithField("context", ctx).Error(err.Error())
	}
}
