# Build Stage
FROM golang:1.21-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/assessment-tax .

# Final Stage
FROM alpine:3.19

WORKDIR /app

COPY --from=build-base /app/out/assessment-tax /app/assessment-tax

# Expose DATABASE_URL environment variable
# ARG DATABASE_URL
# ENV DATABASE_URL=$DATABASE_URL

CMD ["/app/assessment-tax"]