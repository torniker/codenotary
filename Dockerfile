FROM golang:1.22 AS build
WORKDIR /app
ADD go.* .
ADD *.go .
ADD immudb immudb
ADD model model
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/codenotary

FROM node:21 AS build-web
WORKDIR /app
COPY web .
RUN npm install
RUN npm run build

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/codenotary /app/codenotary
COPY --from=build-web /app/dist /app/web/dist
EXPOSE 5656
ENTRYPOINT ["/app/codenotary"]
