FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

# - CGO_ENABLED=0: Disable Cgo to produce a statically-linked binary without external C dependencies.
# -o /app/stock-picker: Specifies the output file path and name.
# -ldflags="-w -s": Strips debugging information, reducing the binary size.
#   -w: Omit the DWARF symbol table.
#   -s: Omit the symbol table and debug information.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/stock-picker /app/cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "/app/stock-picker" ]