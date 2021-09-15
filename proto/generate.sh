
#USER
protoc proto/userpb/user.proto --go_out=plugins=grpc:.
protoc proto/userpb/user.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.
