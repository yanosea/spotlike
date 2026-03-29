FROM golang:1.25
WORKDIR /spotlike
COPY . .
RUN go mod download
CMD ["make", "test.local"]
