module github.com/SandyWalsh/runner/server

go 1.19

replace github.com/SandyWalsh/runner/sprinter => ../sprinter

replace github.com/SandyWalsh/runner/sprinter/proto => ../sprinter/proto

require (
	github.com/SandyWalsh/runner/sprinter v0.0.0-00010101000000-000000000000
	github.com/SandyWalsh/runner/sprinter/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.49.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/net v0.0.0-20220907135653-1e95f45603a7 // indirect
	golang.org/x/sys v0.0.0-20220907062415-87db552b00fd // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220902135211-223410557253 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
