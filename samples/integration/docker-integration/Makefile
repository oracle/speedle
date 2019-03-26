.PHONY: all binary install clean uninstall

LIBDIR=${DESTDIR}/lib/systemd/system
BINDIR=${DESTDIR}/usr/lib/docker

all: build

build:
	go build -o ${HOME}/go/bin/speedle-docker-plugin *.go

install: build
	install -m 644 systemd/speedle-docker-plugin.service ${LIBDIR}
	install -m 644 systemd/speedle-docker-plugin.socket ${LIBDIR}
	install -m 755 ${GOPATH}/bin/speedle-docker-plugin ${BINDIR}

clean:
	rm -f ${GOROOT}/bin/casbin-authz-plugin

uninstall:
	rm -f ${LIBDIR}/speedle-docker-plugin.service
	rm -f ${LIBDIR}/speedle-docker-plugin.socket
	rm -f ${BINDIR}/casbin-authz-plugin
