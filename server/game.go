package server

type Game struct {
	Rooms map[string]*Room
}

func (game *Game) HandleEvent(event *Event) {
	switch event.Type {

	}
}
