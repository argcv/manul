package config

func GetFsWorkdir() string {
	return getStringOrDefault(KeyFsWorkdir, "/tmp/manul/workdir")
}
