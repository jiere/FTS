FROM golang:1.13-alpine as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

# Set up the PROJ env variable
RUN mkdir /project
ENV PROJ /project
COPY . $PROJ

# Set WORKDIR
WORKDIR $PROJ

RUN go mod init fts.local && \
    go get -u github.com/swaggo/swag/cmd/swag && \
    swag init && \
    go build -o app && \ 
    mkdir -p /bin/cfg && \    
    cp $PROJ/app /bin && \
    cp $PROJ/cfg/db.ini /bin/cfg/

# Run unit tests and collect coverage data
#ENV CGO_ENABLE 0
#RUN  ./ut.sh

FROM alpine as prod

RUN mkdir -p /fts/bin
WORKDIR /fts/bin

COPY  --from=build /bin /fts/bin/

RUN chmod +x /fts/bin/app
ENTRYPOINT ["/fts/bin/app"]

EXPOSE 8080/tcp
