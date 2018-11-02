package config

const (
	// RPC Section
	KeyRpcBind                   = "rpc.bind"
	KeyRpcHost                   = "rpc.host"
	KeyRpcPort                   = "rpc.port"
	KeyRpcOptionMaxRecvMsgSizeMB = "rpc.option.max_recv_msg_size_mb"
	KeyRpcOptionMaxSendMsgSizeMB = "rpc.option.max_send_msg_size_mb"
	KeyHttpBind                  = "rpc.http.bind"
	KeyHttpHost                  = "rpc.http.host"
	KeyHttpPort                  = "rpc.http.port"

	// Database
	KeyDBMongoAddrs         = "db.mongo.addrs"
	KeyDBMongoPerformAuth   = "db.mongo.with_auth"
	KeyDBMongoAuthDatabase  = "db.mongo.auth.db"
	KeyDBMongoAuthUser      = "db.mongo.auth.user"
	KeyDBMongoAuthPass      = "db.mongo.auth.pass"
	KeyDBMongoAuthMechanism = "db.mongo.auth.mechanism"
	KeyDBMongoTimeoutSec    = "db.mongo.timeout_sec"
	KeyDBMongoDatabase      = "db.mongo.db"

	// Client
	KeyClientUserName   = "client.user.name"
	KeyClientUserSecret = "client.user.secret"
	//KeyClientProject    = "client.project"

	KeyFsWorkdir = "fs.workdir"

	KeyMailSMTPHost               = "mail.smtp.host"
	KeyMailSMTPPort               = "mail.smtp.port"
	KeyMailSMTPUserSender         = "mail.smtp.sender" // Sender Name, Like: Sunlab Team
	KeyMailSMTPPerformAuth        = "mail.smtp.with_auth"
	KeyMailSMTPUserName           = "mail.smtp.auth.user"
	KeyMailSMTPPassword           = "mail.smtp.auth.pass"
	KeyMailSMTPUserDefaultFrom    = "mail.smtp.default_from"
	KeyMailSMTPInsecureSkipVerify = "mail.smtp.insecure_skip_verify"
)
