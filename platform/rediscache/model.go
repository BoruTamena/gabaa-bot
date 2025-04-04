package rediscache

type RedisClient struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(addr, password string, db int) RedisClient {
	return RedisClient{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}
