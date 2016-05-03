FROM gliderlabs/alpine
RUN adduser -D app
USER app
COPY snakepit-seed /usr/bin
COPY swagger.json /home/app/
COPY config.yaml /home/app/
EXPOSE 3000
ENTRYPOINT ["snakepit-seed"]
