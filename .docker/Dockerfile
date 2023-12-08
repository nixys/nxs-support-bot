FROM golang:1.21.1-bullseye as build-env

WORKDIR /usr/src/app
ADD ./ /usr/src/app

RUN go build -o /app

FROM gcr.io/distroless/base

COPY --from=build-env /app /

COPY ./localization/ /localization/

CMD ["/app", "-c", "/app.conf"]
