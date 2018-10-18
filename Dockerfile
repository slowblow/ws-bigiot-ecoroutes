FROM golang:1.10.2-alpine3.7 as builder
COPY . /go/src/ws-bigiot-ecoroutes
WORKDIR /go/src/ws-bigiot-ecoroutes
RUN apk --no-cache add git openssh-client glide && \
    glide cc && glide update && go build

FROM openjdk:8-jre-alpine as final
RUN apk --no-cache add ca-certificates
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Europe/Madrid /etc/localtime
WORKDIR /root/
COPY --from=builder /go/src/ws-bigiot-ecoroutes/conf ./conf
COPY --from=builder /go/src/ws-bigiot-ecoroutes/ws-bigiot-ecoroutes ./ws-bigiot-ecoroutes
COPY --from=builder /go/src/ws-bigiot-ecoroutes/java-example-consumer.jar ./java-example-consumer.jar
CMD env PORT=8080 ./ws-bigiot-ecoroutes
EXPOSE 8080
