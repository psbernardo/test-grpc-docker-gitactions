# test-grpc-docker-gitactions

Download file base on your OS https://github.com/protocolbuffers/protobuf/releases/tag/v3.12.3

COPY the protoc to your gopath bin directory
2.1 inside include folder, copy google folder and paste in gopath bin directory

Add your [gopath]/bin to path environment variable 
3.1 test in terminal $ protoc

go get github.com/googleapis/googleapis/tree/master/google/type
4.1 copy type dir to [gopath]/bin/google/ from [gopath]/pkg/mod/github.com/googleapis/googleapis@[version-hash]/google/type

go get github.com/golang/protobuf/protoc-gen-go@v1.4.2
Test in terminal $ protoc-gen-go

Download https://github.com/grpc/grpc-web/releases/tag/1.2.0

Rename the download file to 'protoc-gen-grpc-web', copy it to [gopath]/bin
7.1 test in terminal $ protoc-gen-grpc-web

NOTE:
This guide assumes you already installed GO. GO is not required in protobuf, if you are using a different programming language you can change "gopath bin folder" to any folder that can be accessed by PATH environment variable.

test note