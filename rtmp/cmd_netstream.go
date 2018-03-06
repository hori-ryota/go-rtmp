package rtmp

//go:generate go run $DEFDIR/go/cmd/genamf/genamf.go -package rtmp -toml $DEFDIR/net_stream_command/on_status.toml -toml $DEFDIR/net_stream_command/play.toml -toml $DEFDIR/net_stream_command/play2.toml -toml $DEFDIR/net_stream_command/delete_stream.toml -toml $DEFDIR/net_stream_command/close_stream.toml -toml $DEFDIR/net_stream_command/receive_audio.toml -toml $DEFDIR/net_stream_command/receive_video.toml -toml $DEFDIR/net_stream_command/publish.toml -toml $DEFDIR/net_stream_command/seek.toml -toml $DEFDIR/net_stream_command/pause.toml -o cmd_netstream_gen.go

//go:generate go run $DEFDIR/go/cmd/genhandler/genhandler.go -package rtmp -name NetStreamCommandHandler -toml $DEFDIR/net_stream_command/on_status.toml -toml $DEFDIR/net_stream_command/play.toml -toml $DEFDIR/net_stream_command/play2.toml -toml $DEFDIR/net_stream_command/delete_stream.toml -toml $DEFDIR/net_stream_command/close_stream.toml -toml $DEFDIR/net_stream_command/receive_audio.toml -toml $DEFDIR/net_stream_command/receive_video.toml -toml $DEFDIR/net_stream_command/publish.toml -toml $DEFDIR/net_stream_command/seek.toml -toml $DEFDIR/net_stream_command/pause.toml -o cmd_netstream_handler_gen.go

//go:generate go run $DEFDIR/go/cmd/genexector/genexecutor.go -package rtmp -name NetStreamCommander -toml $DEFDIR/net_stream_command/on_status.toml -toml $DEFDIR/net_stream_command/play.toml -toml $DEFDIR/net_stream_command/play2.toml -toml $DEFDIR/net_stream_command/delete_stream.toml -toml $DEFDIR/net_stream_command/close_stream.toml -toml $DEFDIR/net_stream_command/receive_audio.toml -toml $DEFDIR/net_stream_command/receive_video.toml -toml $DEFDIR/net_stream_command/publish.toml -toml $DEFDIR/net_stream_command/seek.toml -toml $DEFDIR/net_stream_command/pause.toml -o netstream_commander_gen.go

func (v PublishingType) String() string {
	return string(v)
}
