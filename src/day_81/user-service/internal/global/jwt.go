package global

import "user-service/internal/middlewares/apisixtoken"

var JwtGenerator *apisixtoken.JWTGenerator

func initJwt() {
	secret := Config.Sub("jwt").GetString("secret")
	apisixConsumerKey := Config.Sub("jwt").GetString("key")
	timeout := Config.Sub("jwt").GetDuration("timeout")
	JwtGenerator = apisixtoken.NewJWTGenerator(secret, apisixConsumerKey, timeout)
}
