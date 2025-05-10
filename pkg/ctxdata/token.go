/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package ctxdata

import "github.com/golang-jwt/jwt"

const Identify = "imooc.com"

func GetJwtToken(secretKey string, iat, seconds int64, uid string, deveicetype string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid
	claims["device"] = deveicetype

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
