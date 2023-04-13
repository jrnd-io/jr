package types

//go:generate $GOPATH/bin/gogen-avro -package types . NetDevice.avsc
//go:generate $GOPATH/bin/gogen-avro -package types . User.avsc

//go:generate register . NetDevice.avsc net-device
//go:generate register . User.avsc	user
