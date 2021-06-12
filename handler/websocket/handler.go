package websocket

import "github.com/nentenpizza/werewolves/wserver"

type Service interface {
	REGISTER(s *wserver.Server)
}

// Initialize creates server, registers services and returns server
func Initialize(config wserver.Settings, services ...Service) *wserver.Server {
	server := wserver.NewServer(config)
	for _, v := range services {
		if v != nil {
			v.REGISTER(server)
		}
	}
	return server
}
