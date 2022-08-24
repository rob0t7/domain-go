package rest

import (
	"net/http"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/memrepository"
)

type RESTServer struct {
	http.Server
	companyService *app.CompanyService
}

func New() *RESTServer {
	var server RESTServer
	server.Addr = ":8080"

	server.companyService = app.NewCompanyService(memrepository.New()) // TODO: This should be using DI
	server.registerHandlers()

	return &server
}

func (s *RESTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

func (s *RESTServer) registerHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/companies", s.CompanyCollectionHandler)
	mux.HandleFunc("/companies/", s.CompanyInstanceHandler)
	s.Handler = mux
}
