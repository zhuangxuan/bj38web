package redis

const tokenMap = "bj38web_tokenMap"

// HsetUsernameToken redis中存储username对应token
func HsetUsernameToken(username string, token string) (res bool, err error) {
	res, err = client.HSet(tokenMap, username, token).Result()

	if err != nil {
		return false, err
	}
	return true, nil
}

// HgetUsernameToken redis中获取username对应的token
func HgetUsernameToken(username string) (res string, err error) {
	res, err = client.HGet(tokenMap, username).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
