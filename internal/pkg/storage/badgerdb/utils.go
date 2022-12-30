package badgerdb

func getUserKey(userID string) []byte {
	return []byte(prefixUserKey + userID)
}
