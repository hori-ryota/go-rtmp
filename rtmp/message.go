package rtmp

//go:generate go run $DEFDIR/go/cmd/genstruct/genstruct.go -package rtmp -toml $DEFDIR/message/message.toml -o message_gen.go
//go:generate stringer -type MessageTypeID -trimprefix MessageTypeID -output message_typeid_string_gen.go
