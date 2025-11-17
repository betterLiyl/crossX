package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

func VerifyHMACSHA256(body map[string]any, secret string, sigHex string) error {
	if secret == "" {
		return nil
	}
	keys := make([]string, 0, len(body))
	for k := range body {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	mac := hmac.New(sha256.New, []byte(secret))
	for _, k := range keys {
		v := body[k]
		mac.Write([]byte(k))
		mac.Write([]byte("="))
		mac.Write([]byte(fmtAny(v)))
		mac.Write([]byte("&"))
	}
	sum := mac.Sum(nil)
	got, err := hex.DecodeString(sigHex)
	if err != nil {
		return err
	}
	if !hmac.Equal(sum, got) {
		return errors.New("bad_sig")
	}
	return nil
}

func fmtAny(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return trimZeros(x)
	case int64:
		return itoa(x)
	case int:
		return itoa(int64(x))
	default:
		return ""
	}
}

func itoa(n int64) string {
	b := []byte{}
	if n == 0 {
		return "0"
	}
	for n > 0 {
		b = append(b, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func trimZeros(f float64) string {
	s := fmt.Sprintf("%.8f", f)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
