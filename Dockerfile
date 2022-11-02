FROM golang:1.17.8-alpine3.15

# set environment variables
ENV USER_NAME=quizcat
ENV HOME=/home/$USER_NAME
ENV APP_NAME=quizcat
ENV APP_HOME=$HOME/$APP_NAME

RUN mkdir -p $APP_HOME

# set working directory
WORKDIR $APP_HOME

# install some packages
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add gcc g++ make git libpq ca-certificates

# install go dependencies
COPY ./go.mod .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

# copy sourse code
COPY . .

# build app
RUN GOOS=linux go build -ldflags="-s -w" -o $APP_NAME

# change user
#RUN addgroup -S $USER_NAME && adduser -S $USER_NAME -G $USER_NAME
#RUN chown -R $USER_NAME:$USER_NAME $APP_HOME $REPOS_DIR
#USER $USER_NAME

# run app
ENTRYPOINT $APP_HOME/$APP_NAME