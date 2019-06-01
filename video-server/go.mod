module video-server

go 1.12

require (
	github.com/gorilla/mux v1.7.2
	github.com/lib/pq v1.1.1
	github.com/satori/go.uuid v1.2.0
)

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
