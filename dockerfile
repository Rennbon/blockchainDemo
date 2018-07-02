
FROM golang:1.10

WORKDIR /app
ADD myapp /app/

ENTRYPOINT ["./myapp"]