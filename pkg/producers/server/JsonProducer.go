package server

import (
	"github.com/ugol/jr/pkg/tpl"
	"log"
	"net/http"
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
			_, err = (respWriter).Write([]byte(","))
			if err != nil {
				log.Println(err.Error())
			}
		}
		_, err := (respWriter).Write(value)
		if err != nil {
			log.Println(err.Error())
		}
	}

}
