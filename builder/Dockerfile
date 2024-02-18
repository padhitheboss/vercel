FROM node:latest
RUN apt update
RUN mkdir /app
ADD bin/builder /app/builder
WORKDIR /app
CMD ["/app/builder"]