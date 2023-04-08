package actions

import (
	// "fmt"
	// "github.com/golang-jwt/jwt/v4"
)

func VerifyToken(token string) (res string, err error) {

	// //Parse token
	// parse_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte("secret"), nil
	// })
	
	// if err != nil {
	// 	return "", err
	// }

	// //Return claims
	// claims := parse_token.Claims.(jwt.MapClaims)
	// return claims["phone"].(string), nil
	return "", nil
}
