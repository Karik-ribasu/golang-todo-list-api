FROM golang:1.20-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/todo-api .

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /out/todo-api /app/todo-api
USER nonroot:nonroot
EXPOSE 5000
ENTRYPOINT ["/app/todo-api"]
