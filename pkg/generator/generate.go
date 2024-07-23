package generator

/// avrogen CLI doesn't support multiple files at the moment so sticking to gogen-avro
///go:generate $GOPATH/bin/avrogen -o ../types/*.go -pkg types ../types/*.avsc
//go:generate $GOPATH/bin/gogen-avro -package types ../types ../types/*.avsc
//go:generate go run generateRegistry.go
