package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type captchaPayload struct {
	A   int   `json:"a"`
	B   int   `json:"b"`
	Exp int64 `json:"exp"`
}

func randomInt(min, max int) int {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(nBig.Int64()) + min
}

func signCaptcha(secret []byte, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func makeCaptchaToken(secret []byte, a, b int, ttl time.Duration) (question string, token string) {
	p := captchaPayload{
		A:   a,
		B:   b,
		Exp: time.Now().Add(ttl).Unix(),
	}
	raw, _ := json.Marshal(p)
	sig := signCaptcha(secret, raw)
	token = base64.RawURLEncoding.EncodeToString(raw) + "." + sig
	question = "¿Cuánto es " + itoa(a) + " + " + itoa(b) + "?"
	return
}

func verifyCaptchaToken(secret []byte, token string) (a, b int, ok bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return 0, 0, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, 0, false
	}
	gotSig := parts[1]
	wantSig := signCaptcha(secret, raw)
	if !hmac.Equal([]byte(gotSig), []byte(wantSig)) {
		return 0, 0, false
	}

	var p captchaPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return 0, 0, false
	}
	if time.Now().Unix() > p.Exp {
		return 0, 0, false
	}
	return p.A, p.B, true
}

// mini itoa para no importar strconv aquí
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [32]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + (n % 10))
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func captchaSecretFromEnv(c *gin.Context) []byte {
	sec := c.MustGet("CAPTCHA_SECRET").([]byte)
	if len(sec) < 16 {
		// fallback (no ideal, pero evita crash)
		return []byte("change-me-please-32-bytes-min")
	}
	return sec
}

func NewCaptchaForForm(c *gin.Context) (question string, token string) {
	secret := captchaSecretFromEnv(c)
	a := randomInt(1, 9)
	b := randomInt(1, 9)
	return makeCaptchaToken(secret, a, b, 5*time.Minute)
}

func CheckCaptcha(c *gin.Context) bool {
	token := c.PostForm("captcha_token")
	ans := c.PostForm("captcha_answer")

	a, b, ok := verifyCaptchaToken(captchaSecretFromEnv(c), token)
	if !ok {
		return false
	}

	// parse answer simple
	sum := a + b
	return ans == itoa(sum)
}

func RenderCaptchaError(c *gin.Context, view string, title string, breadcrumbs []Crumb, extra gin.H) {
	q, t := NewCaptchaForForm(c)
	data := gin.H{
		"Title":        title,
		"Breadcrumbs":  breadcrumbs,
		"CaptchaQ":     q,
		"CaptchaToken": t,
		"Error":        "Captcha inválido. Inténtalo otra vez.",
	}
	for k, v := range extra {
		data[k] = v
	}
	c.HTML(http.StatusBadRequest, view, data)
}
