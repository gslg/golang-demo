module redis-demo

go 1.12

require (
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/gorilla/mux v1.7.2
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
