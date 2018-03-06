package rtmp

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package rtmp -toml $DEFDIR/message/protocol_control_message/set_chunk_size.toml -toml $DEFDIR/message/protocol_control_message/abort_message.toml -toml $DEFDIR/message/protocol_control_message/acknowledgement.toml -toml $DEFDIR/message/protocol_control_message/window_acknowledgement_size.toml -toml $DEFDIR/message/protocol_control_message/set_peer_bandwidth.toml -o msg_protocol_control_gen.go

//go:generate go run $DEFDIR/go/cmd/genhandler/genhandler.go -package rtmp -name ProtocolControlEventHandler -toml $DEFDIR/message/protocol_control_message/set_chunk_size.toml -toml $DEFDIR/message/protocol_control_message/abort_message.toml -toml $DEFDIR/message/protocol_control_message/acknowledgement.toml -toml $DEFDIR/message/protocol_control_message/window_acknowledgement_size.toml -toml $DEFDIR/message/protocol_control_message/set_peer_bandwidth.toml -o msg_protocol_control_handler_gen.go

//go:generate go run $DEFDIR/go/cmd/genexector/genexecutor.go -package rtmp -name ProtocolController -toml $DEFDIR/message/protocol_control_message/set_chunk_size.toml -toml $DEFDIR/message/protocol_control_message/abort_message.toml -toml $DEFDIR/message/protocol_control_message/acknowledgement.toml -toml $DEFDIR/message/protocol_control_message/window_acknowledgement_size.toml -toml $DEFDIR/message/protocol_control_message/set_peer_bandwidth.toml -o protocol_controller_gen.go

//go:generate stringer -type BandwidthLimitType -trimprefix BandwidthLimitType -output msg_protocol_control_bandwidth_limit_type_string_gen.go
