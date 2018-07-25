package config

func GetRpcBind() string {
	return getStringOrDefault(KeyRpcBind, "0.0.0.0")
}

func GetRpcHost() string {
	return getStringOrDefault(KeyRpcHost, "127.0.0.1")
}

func GetRpcPort() int {
	return getIntOrDefault(KeyRpcPort, 35000)
}

func GetRpcOptionMaxRecvMsgSizeMB() int {
	// 64MB in default
	return getIntOrDefault(KeyRpcOptionMaxRecvMsgSizeMB, 64)
}

func GetRpcOptionMaxSendMsgSizeMB() int {
	// 64MB in default
	return getIntOrDefault(KeyRpcOptionMaxSendMsgSizeMB, 64)
}

func GetRpcOptionMaxRecvMsgSize() int {
	// 64MB in default
	return GetRpcOptionMaxRecvMsgSizeMB() * 1024 * 1024
}

func GetRpcOptionMaxSendMsgSize() int {
	// 64MB in default
	return GetRpcOptionMaxSendMsgSizeMB() * 1024 * 1024
}

func GetHttpBind() string {
	return getStringOrDefault(KeyHttpBind, "0.0.0.0")
}

func GetHttpHost() string {
	return getStringOrDefault(KeyHttpHost, "127.0.0.1")
}

func GetHttpPort() int {
	return getIntOrDefault(KeyHttpPort, 35100)
}
