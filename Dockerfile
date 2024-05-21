FROM debian:stable-slim

COPY files/ /bin/files/
COPY ./chirpy /bin/chirpy

ENV PORT 8080
CMD ["/bin/chirpy"]