FROM ubuntu:latest
RUN mkdir /app
RUN apt update
ADD bin/request-handler /app/request-handler
WORKDIR /app
CMD ["/app/request-handler"]