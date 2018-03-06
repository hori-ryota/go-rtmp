package rtmp

func GenerateCommonConnInitializer(
	messageHandlerFactories ...func(c Conn) MessageHandler,
) func(c Conn) {
	return func(c Conn) {
		messageHandlerFactories = append(messageHandlerFactories, NewConnControllMessageHandler)
		mhs := make([]MessageHandler, len(messageHandlerFactories))
		for i, f := range messageHandlerFactories {
			mhs[i] = f(c)
		}
		c.AppendMessageHandler(mhs...)
	}
}
