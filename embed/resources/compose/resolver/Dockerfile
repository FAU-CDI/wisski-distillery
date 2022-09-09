FROM docker.io/library/alpine

COPY wdcli /wdcli
EXPOSE 8888
CMD ["/wdcli","--internal-in-docker","--config","${CONFIG_PATH}","resolver_server","--bind","0.0.0.0:8888"]