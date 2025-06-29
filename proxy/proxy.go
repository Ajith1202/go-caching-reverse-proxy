package proxy

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
)

func generateCacheKey(r *http.Request) string {
	sb := strings.Builder{}
	sb.WriteString(r.Method)
	sb.WriteString(":")
	sb.WriteString(r.URL.Path)
	sb.WriteString("?")
	sb.WriteString(r.URL.RawQuery)

	// headers := []string{"Accept", "Accept-Encoding"}
	// for _, h := range headers {
	// 	sb.WriteString("|")
	// 	sb.WriteString(h)
	// 	sb.WriteString("=")
	// 	sb.WriteString(r.Header.Get(h))
	// }

	// Hashing for reducing size
	hash := md5.Sum([]byte(sb.String()))

	return hex.EncodeToString(hash[:])
}
