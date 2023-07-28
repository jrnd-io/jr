package test

import (
	"bytes"
	"github.com/ugol/jr/pkg/tpl"
	"log"
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
			_, err = (respWriter).Write([]byte(","))
			if err != nil {
				log.Println(err.Error())
			}
		}
		_, err := (respWriter).Write(value)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Printf("Test producer must produce to a bytes.Buffer, but was a %T\n", o)
	}

}
