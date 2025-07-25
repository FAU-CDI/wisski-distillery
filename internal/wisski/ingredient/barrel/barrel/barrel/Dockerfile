# ============================
# WissKI Distillery Dockerfile
# ============================

# This file is part of the WissKI Distillery and sets up an image
# to be used for individual WissKIs.

# Start from a base image (configured by the build argument).
ARG BARREL_BASE_IMAGE=docker.io/library/php:8.3-apache-bookworm
FROM $BARREL_BASE_IMAGE

# Setup in /var/www
WORKDIR /var/www

# install and enable the various required php extensions and dropbear ssh server
ARG IIP_SERVER_ENABLED=
ARG PHP_CONFIG_MODE=prod
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        curl \
        default-mysql-client \
        git \
        imagemagick \
        libavif-dev \
        libcurl4-openssl-dev \
        libfreetype6-dev \
        libicu-dev \
        libjpeg62-turbo-dev \
        libpng-dev \
        libssh2-1-dev \
        libwebp-dev \
        libxml2-dev \
        libxpm-dev \
        libxslt-dev \
        libzip-dev \
        openssh-server \
        sudo \
        unzip \
        vim \
        zip \
        $([ "$IIP_SERVER_ENABLED" = "1" ] && echo iipimage-server) \
    && \
    docker-php-source extract && \
    mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini" && \
    pear config-set php_ini "$PHP_INI_DIR/php.ini" && \
    docker-php-ext-configure gd \
        --enable-gd \
        --with-avif \
        --with-webp \
        --with-jpeg \
        --with-xpm \
        --with-freetype \
        --enable-gd-jis-conv \
    && \
    docker-php-ext-install \
        curl \
        gd \
        intl \
        mysqli \
        opcache \
        pdo_mysql \
        soap \
        xml \
        xmlwriter \
        xsl \
        zip \
    && \
    pecl install uploadprogress-2.0.2 && \
    pecl install xmlrpc-1.0.0RC3 && \
    pecl install ssh2-1.4.1 && \
    pecl install apcu-5.1.24 && \
    (([ "$PHP_CONFIG_MODE" = "devel" ] && pecl install xdebug-3.4.3) || true) && \
    docker-php-ext-enable \
        apcu \
        curl \
        gd \
        intl \
        mysqli \
        mysqli \
        opcache \
        pdo_mysql \
        soap \
        ssh2 \
        uploadprogress \
        xml \
        xmlrpc \
        xmlwriter \
        xsl \
        zip \
    && \
    docker-php-source delete

# enable the apache rewrite mod
RUN a2enmod rewrite headers && \
    (([ "$PHP_CONFIG_MODE" = "devel" ] && touch /var/log/xdebug.log && chown www-data:www-data /var/log/xdebug.log) || true)

# Install composer.
ARG COMPOSER_VERSION=2.8.9
RUN curl -sS https://getcomposer.org/installer | php -- --version=$COMPOSER_VERSION && \
    mv composer.phar /usr/local/bin/composer

# Add it to the path
ENV PATH "/usr/local/bin:/var/www/data/project/vendor/bin:$PATH"

# Configure PHP
ADD php.ini.d/wisski.ini /usr/local/etc/php/conf.d/01_wisski.ini

# Configure the php config with prod or devel
ADD php.ini.d/php-$PHP_CONFIG_MODE.ini /usr/local/etc/php/conf.d/02_mode.ini

ARG CONTENT_SECURITY_POLICY=
ENV CONTENT_SECURITY_POLICY=${CONTENT_SECURITY_POLICY}

# Configure Apache.

# first remove the default configuration
RUN rm /etc/apache2/sites-available/*.conf && \
    rm /etc/apache2/sites-enabled/*.conf

# Then add all the configs
ADD apache.d/conf/ports.conf /etc/apache2/ports.conf
ADD apache.d/sites-available/wisski.conf /etc/apache2/sites-available/wisski.conf

# this file is technically only needed for iipserv is enabled
# but we add it in either case, as it's not enabled.
ADD apache.d/mods-available/iipsrv.conf /etc/apache2/mods-available/iipsrv.conf

# enable it
RUN a2ensite wisski

# volumes for all the things
VOLUME /var/www/.composer
VOLUME /var/www/data
VOLUME /ssh/hostkeys/

# Add and configure the entrypoint
ADD etc/entrypoint.sh /entrypoint.sh
ADD etc/profile_path.sh /etc/profile.d/path.sh
ADD ssh/ /ssh/

RUN chmod 500 /ssh/keys.sh && \
    chmod 500 /ssh/start.sh && \
    chmod 444 /etc/profile.d/path.sh && \
    chmod 400 /entrypoint.sh && \
    chsh www-data --shell /bin/bash

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]

CMD ["apache2-foreground"]
EXPOSE 8080