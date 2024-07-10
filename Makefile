BUILD ?= build
BINARY ?= glone
GO ?= go
GLONE_VERSION ?= v1.0.0
LDFLAGS = -ldflags "-X main.Version=${GLONE_VERSION}($(shell git rev-parse --short HEAD))"
glone:
	@mkdir -p ${BUILD}
	${GO} build -o ${BUILD}/${BINARY} ${LDFLAGS} .

clean:
	rm -rf build/*

.PHONY: glone clean
