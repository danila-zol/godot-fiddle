FROM golang:1.24.2 AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

ARG MAXPROCS=1
COPY . .
RUN go test -p 1 ./...

# Creates a staticaly-linked binary
RUN CGO_ENABLED=0 go build -v -o ./bin/ -gcflags=GOMAXPROCS=$MAXPROCS ./...

# Docs generation using swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN export PATH=$(go env GOPATH)/bin:$PATH
RUN swag init --dir ./cmd/api/,./internal/delivery/http/v1,./internal/domain/models --parseInternal

FROM scratch

COPY --from=builder /usr/src/app/docs/ ./
COPY --from=builder /usr/src/app/bin/api ./
COPY --from=builder /usr/src/app/internal/enforcer/psqlCasbinClient/rbac_model.conf ./internal/enforcer/psqlCasbinClient/rbac_model.conf

# For S3 operations
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "/api" ]
