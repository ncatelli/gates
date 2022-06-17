
ARG BASEIMG="alpine:3.16"
ARG BUILDIMG="golang:1.17-alpine"
FROM $BUILDIMG as builder

ARG APP_NAME="gates"
ENV GOPATH=""

COPY . /go/

RUN cd /go && \
	go build -o /${APP_NAME}

FROM $BASEIMG
LABEL maintainer="Nate Catelli <ncatelli@packetfire.org>"
LABEL description="Container for gates"

ARG SERVICE_USER="service"
ARG APP_NAME="gates"

RUN addgroup ${SERVICE_USER} && \
	adduser -D -G ${SERVICE_USER} ${SERVICE_USER}

COPY --from=builder /${APP_NAME} /opt/${APP_NAME}/bin/${APP_NAME}

RUN chown -R ${SERVICE_USER}:${SERVICE_USER} /opt/${APP_NAME}/bin/${APP_NAME} && \
	chmod +x /opt/${APP_NAME}/bin/${APP_NAME}

RUN apk --no-cache add curl

WORKDIR "/opt/${APP_NAME}/"
USER ${SERVICE_USER}

ENTRYPOINT [ "/opt/gates/bin/gates" ]
CMD [ "-h" ]