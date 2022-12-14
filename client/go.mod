module github.com/SandyWalsh/runner/client

go 1.19

replace github.com/SandyWalsh/runner/sprinter/proto => ../sprinter/proto

require (
	github.com/SandyWalsh/runner/sprinter/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.49.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220909162455-aba9fc2a8ff2 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220909194730-69f6226f97e5 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
