FROM alpine:latest
WORKDIR /root/
COPY cloud_firewall .
EXPOSE 8080
CMD ["./cloud_firewall"]