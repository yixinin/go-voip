FROM centos:latest

Copy voip voip
COPY config/docker.yaml config/app.yaml

CMD [ "./voip"]