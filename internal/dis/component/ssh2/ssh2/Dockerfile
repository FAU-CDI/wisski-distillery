FROM docker.io/library/docker:28.5-cli

COPY wdcli /wdcli
EXPOSE 2222
CMD ["/wdcli","--internal-in-docker","ssh","--private-key-path", "/data/", "--bind","0.0.0.0:2222"]