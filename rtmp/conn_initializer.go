package rtmp

func GenerateCommonConnInitializer() func(c Conn) {
	return func(c Conn) {
		c.AddMessageHandler("ControlMessageHandler", NewControlMessageHandler(c))
	}
}
