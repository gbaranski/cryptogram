package chat

// RoomBufSize is the number of incoming messages to buffer for each topic.
const RoomBufSize = 128

// DefaultTopicName is default topic name
const DefaultTopicName = "general-1"

// GetTopicName retreives topic name from room name
func GetTopicName(roomName *string) string {
	return "cryptogram:" + *roomName
}
