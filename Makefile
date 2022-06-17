APP_NAME="gates"
IMGNAME="ncatelli/${APP_NAME}"
PKG="github.com/ncatelli/${APP_NAME}"

build: 
	go build

build-docker: fmt test
	docker build -t ${IMGNAME}:latest .

test:
	go test -race -cover ./...

benchmark:
	go test -benchmem -bench . ./...

fmt:
	test -z $(shell go fmt ./...)

clean-docker:
	@type docker >/dev/null 2>&1 && \
	docker rmi -f ${IMGNAME}:latest || \
	true

clean: clean-docker
	@rm -f ${APP_NAME} || true

lint:
	golint -set_exit_status ./...
