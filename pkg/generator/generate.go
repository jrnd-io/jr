package generator

//go:generate $GOPATH/bin/gogen-avro -package types ../types ../types/*.avsc
//go:generate go run generateRegistry.go
