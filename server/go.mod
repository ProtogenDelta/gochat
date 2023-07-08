module github.com/protogendelta/gochat/server

go 1.20

replace github.com/protogendelta/gochat/lib => ../lib

require (
	github.com/protogendelta/gochat/lib v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.31.0
)
