BUILD ?= build
BINARY ?= glone
GO ?= go
glone:
	@mkdir -p ${BUILD}
	${GO} build -o ${BUILD}/${BINARY} .

clean:
	rm -rf build/*

.PHONY: glone clean
