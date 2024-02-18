FROM ubuntu:latest
RUN mkdir /app
ADD bin/kube-job-scheduler /app/kube-job-scheduler
WORKDIR /app
CMD ["/app/kube-job-scheduler"]