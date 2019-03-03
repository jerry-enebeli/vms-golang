package mws

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	httplib "bitbucket.org/eviconlabs/platphom_backend/delivery/http/libs/http"

	"bitbucket.org/eviconlabs/platphom_backend/config/responses"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
)

type exception struct {
	Message string `json:"message"`
}

//Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

//ChainMiddleware chains multiply handlers
func ChainMiddleware(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

//AuthorizationSingle with jwt
func AuthorizationSingle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					resp := &responses.GeneralResponse{Success: false, Message: "token error", Error: error.Error()}
					httplib.Response(w, resp)
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					w.Write([]byte("Invalid authorization token"))

				}
			}
		} else {
			w.Write([]byte("An authorization header is required"))

		}
	})
}

//AuthorizationChain with jwt
func AuthorizationChain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					resp := &responses.GeneralResponse{Success: false, Message: "token error", Error: error.Error()}
					httplib.Response(w, resp)
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next.ServeHTTP(w, req)
				} else {
					w.Write([]byte("Invalid authorization token"))

				}
			}
		} else {
			w.Write([]byte("An authorization header is required"))

		}
	})
}

// //AccessLogToFile saves sever logs to a specifiied file (access.log)
// func AccessLogToFile(r http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		f, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 		if err != nil {
// 			log.Panic("Access log: ", err)
// 		}

// 		logger := handlers.CombinedLoggingHandler(f, r)

// 		logger.ServeHTTP(w, req)
// 	})
// }

//AccessLogToConsole prints sever logs to the terminal
func AccessLogToConsole(r http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		logger := handlers.CombinedLoggingHandler(os.Stdout, r)

		logger.ServeHTTP(w, req)

	})
}