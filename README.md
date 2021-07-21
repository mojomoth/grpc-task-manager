# GRPC-TASK-MANAGER

## grpc setup

1. 
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

2. 
```
protoc -I=. \
	    --go_out . --go_opt paths=source_relative \
	    --go-grpc_out . --go-grpc_opt paths=source_relative \
	    protos/v1/task/task.proto
```


