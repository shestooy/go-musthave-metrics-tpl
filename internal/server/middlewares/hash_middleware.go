package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
)

type responseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func hash(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return h.Sum(nil)
}

func Hash(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if key == "" {
				return next(c)
			}

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				c.Error(err)
			}
			c.Request().Body = io.NopCloser(strings.NewReader(string(body)))
			bodyHash := hash(body, key)

			reqBodyHash := c.Request().Header.Get("HashSHA256")
			resHash, err := hex.DecodeString(reqBodyHash)
			if err != nil {
				c.Error(err)
			}

			if !hmac.Equal(bodyHash, resHash) {
				return c.JSON(http.StatusBadRequest, map[string]string{"err": "the hash checksum did not match"})
			}

			resBody := new(strings.Builder)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			c.Response().Writer = &responseWriter{Writer: mw, ResponseWriter: c.Response().Writer}

			err = next(c)
			if err != nil {
				return err
			}

			responseBytes := []byte(resBody.String())
			resHash = hash(responseBytes, key)
			c.Response().Header().Set("HashSHA256", hex.EncodeToString(resHash))

			return err
		}
	}
}
