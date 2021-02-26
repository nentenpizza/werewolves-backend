package websocket

import "github.com/mitchellh/mapstructure"

func (s *handler) mapToEvent(event interface{}, js interface{}) error {
	err := mapstructure.Decode(js, event)
	if err != nil{
		return err
	}
	return nil
}