package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
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

			reqBodyHash := c.Request().Header.Get("HashSHA256")
			if reqBodyHash == "" {
				return next(c)
			}

			resHash, err := hex.DecodeString(reqBodyHash)
			if err != nil {
				return err
			}

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}

			c.Request().Body = io.NopCloser(io.MultiReader(bytes.NewReader(body)))
			bodyHash := hash(body, key)

			if !hmac.Equal(bodyHash, resHash) {
				return c.JSON(http.StatusBadRequest, map[string]string{"err": "the hash checksum did not match"})
			}

			res := &responseWriter{
				ResponseWriter: c.Response().Writer,
				Body:           &bytes.Buffer{},
			}
			c.Response().Writer = res
			err = next(c)
			if err != nil {
				return err
			}

			responseBytes := res.Body.Bytes()
			resHash = hash(responseBytes, key)
			c.Response().Header().Set("HashSHA256", hex.EncodeToString(resHash))

			return err
		}
	}
}
