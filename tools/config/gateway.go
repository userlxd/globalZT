package config

type Gateway struct {
	LOG  LOG
	AUTH AUTH
	GRPC GRPC
}

type GRPC struct {
	PORT uint
}

type AUTH struct {
	SigningKey  string `default:"0987654321"`
	ExpiresTime int64  `default:"600"`
	BufferTime  int64  `default:"600"`
}

func (c *Gateway) Load() {
	env.ignorePrefix = true
	env.Fill(c)
}
