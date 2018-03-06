package rtmp

import "context"

func SetChunkStreamIDToContext(ctx context.Context, chunkStreamID uint32) context.Context {
	return context.WithValue(ctx, "chunkStreamID", chunkStreamID)
}

func ChunkStreamIDFromContext(ctx context.Context) uint32 {
	return ctx.Value("chunkStreamID").(uint32)
}

func SetMessageStreamIDToContext(ctx context.Context, messageStreamID uint32) context.Context {
	return context.WithValue(ctx, "messageStreamID", messageStreamID)
}

func MessageStreamIDFromContext(ctx context.Context) uint32 {
	return ctx.Value("messageStreamID").(uint32)
}
