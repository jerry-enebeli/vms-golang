package routes

import (
	"net/http"
	"vms/config/responses"
	httplib "vms/libs/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//Router for all routes
func Router() *mux.Router {
	route := mux.NewRouter()

	//BASE ROUTE
	route.HandleFunc("/v1", func(res http.ResponseWriter, req *http.Request) {
		resp := &responses.GeneralResponse{Success: true, Message: "platphom  server running....", Data: "platphom SERVER v1.0"}
		httplib.Response(res, resp)
	})

	//	APPLY MIDDLEWARES
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	})

	route.Use(mws.AccessLogToConsole)
	c.Handler(route)

	mwsWithAuth := mws.ChainMiddleware(mws.AuthorizationSingle)

	//*****************
	// USER ROUTES
	//*****************

	return route
}
