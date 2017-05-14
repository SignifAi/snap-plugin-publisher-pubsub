go:
				glide up
				sed -i 's:grpc.SupportPackageIsVersion3:grpc.SupportPackageIsVersion4:g' vendor/github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc/plugin.pb.go
				go build
all:
				go
clean:
				rm -rf snap-plugin-publisher-pubsub
test:
				gcloud beta emulators pubsub start --host-port=localhost:8321 &
				go test -v $$(glide novendor)
				# kinda dirty
				disown %1
				kill -9 $(shell ps aux |grep gcloud | awk {'print $$2'})
