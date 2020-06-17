# This Dockerfile contains instructions to compile and run GraphDB inside a Docker container. 
# It is roughly based on https://github.com/Ontotext-AD/graphdb-docker/blob/master/free-edition/Dockerfile
# but has been modified for performance and security. 

# This image is intended to be built like:
#   docker build --build-arg graphdb_src=graphdb.zip .

# We first make a base image to base further builds on. 
# We don't use alpine here, as that uses significantly slower musl instead of glibc. 
FROM adoptopenjdk/openjdk11:debian-slim as base

# Create a user called graphdb
RUN useradd -ms /bin/bash graphdb

# make a base images, to add the sources to. 
FROM base as sources

# install unzip
RUN apt-get update && apt-get install -y unzip

# add the source file (by default graphdb.zip) to the image
ARG src=graphdb.zip
ADD ${src} /graphdb.zip

# unpack it into a temporary directory
RUN unzip "$src" -d "/unpack/"

# Move it into /opt/graphdb, and chown it to graphdb
RUN mv "/unpack"/* /opt/graphdb
RUN chown -R graphdb:graphdb /opt/graphdb

# finally make an image that will run
FROM base as final

# add the entrypoint script
ADD entrypoint.sh /entrypoint.sh

# copy over the sources
COPY --from=sources /opt/graphdb /opt/graphdb

# set environment variables for graphdb_home and path
ENV GRAPHDB_HOME=/opt/graphdb
ENV PATH=$GRAPHDB_HOME/bin:$PATH

# expose a port
EXPOSE 7200

# setup a healthcheck, that checks if the server is up. 
RUN apt-get update && apt-get install -y curl
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl --fail 127.0.0.1:7200/rest/repositories || exit 1

# Add volumes for data, work and logs as these might be accessible from the outside. 
# To add your own configuration, manually mount a config file into /opt/graphdb/work
VOLUME /opt/graphdb/data
VOLUME /opt/graphdb/work
VOLUME /opt/graphdb/logs

# setup command and entrypoint
CMD ["-Dgraphdb.home=/opt/graphdb"]
ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]
