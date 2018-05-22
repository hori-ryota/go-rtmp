package rtmp

func GenerateCommonConnInitializer(
	messageHandlerFactories ...func(c Conn) MessageHandler,
) func(c Conn) {
	return func(c Conn) {
		mhfs := append(messageHandlerFactories, func(c Conn) MessageHandler {
			return NewControlMessageHandler(c)
		})
		mhs := make([]MessageHandler, len(mhfs))
		for i, f := range mhfs {
			mhs[i] = f(c)
		}
		c.AppendMessageHandler(mhs...)
	}
}
