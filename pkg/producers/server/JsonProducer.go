package server

import (
	"encoding/json"
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

	data := struct {
		K string
		V string
	}{string(key), string(value)}

	out, err := json.Marshal(c.OutputTpl.ExecuteWith(data))
	if err != nil {
		log.Println(err.Error())
	}
	if o != nil {
		respWriter := o.(*http.ResponseWriter)
		_, err := (*respWriter).Write(out)
		if err != nil {
			log.Println(err.Error())
		}
	}

}
