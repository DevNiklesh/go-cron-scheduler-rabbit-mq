FROM public.ecr.aws/bitnami/golang:1.17 AS cronapp
WORKDIR /go/src/app
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o cronapp .

FROM scratch
COPY --from=cronapp /go/src/app/cronapp /cronapp
CMD ["/cronapp"]