package auth

import(
	"net/http"
	"strings"
	"fmt"
)

func GetAPIKey(headers http.Header) (string, error) {
	rawKey := headers.Get("Authorization")
	token, ok := strings.CutPrefix(rawKey, "ApiKey ")
	if !ok {
		return rawKey, fmt.Errorf("unauthorized token")
	}
	return token, nil
}