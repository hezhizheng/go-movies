FROM scratch

WORKDIR $GOPATH/src/github.com/hezhizheng/go-movies
COPY . $GOPATH/src/github.com/hezhizheng/go-movies

#EXPOSE 8000, 默认直接使用编译好的端口
CMD ["./go_movies_linux_amd64"]
