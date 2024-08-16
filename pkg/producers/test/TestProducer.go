package test

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/ugol/jr/pkg/tpl"
)

type TestProducer struct {
	OutputTpl *tpl.Tpl
}

func (c *TestProducer) Close() error {
	// no need to close
	return nil
}

func (c *TestProducer) Produce(key []byte, value []byte, o interface{}) {

	if o != nil {
		respWriter := o.(*bytes.Buffer)
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
		log.Warn().Interface("o", o).Msg("Test producer must produce to a bytes.Buffer")
	}

}
