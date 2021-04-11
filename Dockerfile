FROM golang:1.16-alpine as build

ENV CGO_ENABLED=0

COPY ./ /go/src/github.com/meyskens/helmet-controller
WORKDIR /go/src/github.com/meyskens/helmet-controller

RUN go build ./

FROM scratch
COPY --from=build /go/src/github.com/meyskens/helmet-controller/helmet-controller /
ENTRYPOINT ["/helmet-controller"]