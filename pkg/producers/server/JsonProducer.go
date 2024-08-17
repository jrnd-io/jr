package server

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ugol/jr/pkg/tpl"
)

type JsonProducer struct {
	OutputTpl *tpl.Tpl
}

func (c *JsonProducer) Close() error {
	// no need to close
	return nil
}

func (c *JsonProducer) Produce(key []byte, value []byte, o interface{}) {

	if o != nil {
		respWriter := o.(http.ResponseWriter)
		if string(key) != "null" {
			_, err := (respWriter).Write(key)
			if err != nil {
				log.Error().Err(err).Msg("Error writing key")
			}
			_, err = (respWriter).Write([]byte(","))
			if err != nil {
				log.Error().Err(err).Msg("Error writing comma")
			}
		}
		_, err := (respWriter).Write(value)
		if err != nil {
			log.Error().Err(err).Msg("Error writing value")
		}
	} else {
		log.Warn().Interface("o", o).Msg("Server producer must produce to a http.ResponseWriter")
	}

}
