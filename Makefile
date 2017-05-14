go:
				glide up
				go build
all:
				go
clean:
				rm -rf snap-plugin-publisher-pubsub
test:
				gcloud beta emulators pubsub start --host-port=localhost:8321 > /dev/null 2>&1 &
				go test -v $$(glide novendor)
				# kinda dirty
				kill -9 $(shell ps aux |grep gcloud | awk {'print $$2'})
