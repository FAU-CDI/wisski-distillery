ARG RDF4J_VERSION=latest
FROM eclipse/rdf4j-workbench:${RDF4J_VERSION}

ENV RDF4J_REPOSITORY="default" \
    RDF4J_REPOSITORY_LABEL="Automatically created Repository"
COPY entrypoint.sh healthcheck.sh /

EXPOSE 8080

VOLUME /var/rdf4j
VOLUME /usr/local/tomcat/logs

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "catalina.sh", "run" ]

COPY ROOT /usr/local/tomcat/webapps/ROOT

HEALTHCHECK --interval=1m --timeout=5s --start-period=30s --start-interval=5s --retries=3 CMD [ "/bin/bash", "/healthcheck.sh" ]
