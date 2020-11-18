module github.com/NOVAPokemon/store

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/gorilla/mux v1.7.4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
)

replace (
	github.com/NOVAPokemon/utils v0.0.62 => ../utils
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2 => ./../../bruno-anjos/archimedesHTTPClient
)
