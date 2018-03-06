package rtmp

//go:generate go run $DEFDIR/go/cmd/genamf/genamf.go -package rtmp -toml $DEFDIR/net_connection_command/connect.toml -toml $DEFDIR/net_connection_command/call.toml -toml $DEFDIR/net_connection_command/close.toml -toml $DEFDIR/net_connection_command/create_stream.toml -o cmd_netconnection_gen.go

//go:generate go run $DEFDIR/go/cmd/genhandler/genhandler.go -package rtmp -name NetConnectionCommandHandler -toml $DEFDIR/net_connection_command/connect.toml -toml $DEFDIR/net_connection_command/call.toml -toml $DEFDIR/net_connection_command/close.toml -toml $DEFDIR/net_connection_command/create_stream.toml -o cmd_netconnection_handler_gen.go

//go:generate go run $DEFDIR/go/cmd/genexector/genexecutor.go -package rtmp -name NetConnectionCommander -toml $DEFDIR/net_connection_command/connect.toml -toml $DEFDIR/net_connection_command/call.toml -toml $DEFDIR/net_connection_command/create_stream.toml -o netconnection_commander_gen.go

//go:generate stringer -type AudioCodecFlag -trimprefix AudioCodecFlag -output cmd_netconnection_audio_codec_flag_string_gen.go
//go:generate stringer -type VideoCodecFlag -trimprefix VideoCodecFlag -output cmd_netconnection_video_codec_flag_string_gen.go
//go:generate stringer -type VideoFunctionFlag -trimprefix VideoFunctionFlag -output cmd_netconnection_video_function_flag_string_gen.go
//go:generate stringer -type EncodingAMFType -trimprefix EncodingAMFType -output cmd_netconnection_encoding_amf_type_string_gen.go
