FROM ubuntu:latest
RUN mkdir /app
RUN apt update && apt install -y git
ADD bin/uploader /app/uploader
WORKDIR /app
CMD ["/app/uploader"]