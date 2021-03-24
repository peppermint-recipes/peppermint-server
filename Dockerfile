FROM golang:1.16-buster AS build
WORKDIR /opt
COPY . .
RUN go build -o /opt/peppermint-server .

FROM scratch AS server
COPY --from=build /opt/peppermint-server /
WORKDIR /opt
CMD ["./peppermint-server"]
