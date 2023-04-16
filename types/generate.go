package types

//go:generate $GOPATH/bin/gogen-avro -package types . NetDevice.avsc
//go:generate $GOPATH/bin/gogen-avro -package types . User.avsc
//go:generate go run generateRegistry.go
