module app

go 1.22

//replace app/config => ./internal/config
//replace app/controllers => ./internal/controllers

require (
	github.com/gorilla/mux v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)
