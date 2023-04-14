package types

//go:generate $GOPATH/bin/gogen-avro -package types . NetDevice.avsc
//go:generate $GOPATH/bin/gogen-avro -package types . User.avsc
//go:generate go run instance.go NetDevice
//go:generate go run instance.go User
