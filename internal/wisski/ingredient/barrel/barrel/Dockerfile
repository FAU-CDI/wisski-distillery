# ============================
# WissKI Distillery Dockerfile
# ============================

# This file is part of the WissKI Distillery and sets up an image
# to be used for individual WissKIs.

# Start from a base image (configured by the build argument).
ARG BARREL_BASE_IMAGE=docker.io/library/php:8.1-apache-bullseye
FROM $BARREL_BASE_IMAGE

# Setup in /var/www
WORKDIR /var/www

# install and enable the various required php extensions and dropbear ssh server
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        curl \
        openssh-server \
        default-mysql-client \
        git \
        imagemagick \
        libcurl4-openssl-dev \
        libfreetype6-dev \
        libicu-dev \
        libjpeg62-turbo-dev \
        libpng-dev \
        libssh2-1-dev \
        libwebp-dev \
        libxml2-dev \
        libxpm-dev \
        sudo \
        unzip \
        vim \
        zip \
    && \
    docker-php-source extract && \
    mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini" && \
    pear config-set php_ini "$PHP_INI_DIR/php.ini" && \
    docker-php-ext-configure gd \
        --enable-gd \
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
    && \
    pecl install xmlrpc-1.0.0RC3 && \
    pecl install ssh2-1.3.1 && \
    pecl install apcu-5.1.21 && \
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
        xml \
        xmlrpc \
    && \
    docker-php-source delete

# enable the apache rewrite mod
RUN a2enmod rewrite headers


# Install composer.
ARG COMPOSER_VERSION=2.6.5
RUN curl -sS https://getcomposer.org/installer | php -- --version=$COMPOSER_VERSION && \
    mv composer.phar /usr/local/bin/composer

# Add it to the path
ENV PATH "/usr/local/bin:/var/www/data/project/vendor/bin:$PATH"

# Configure PHP
ADD php.ini.d/wisski.ini /usr/local/etc/php/conf.d/wisski.ini

# Configure opcache with whatever the user configured
ARG OPCACHE_MODE=prod
ADD php.ini.d/opcache-$OPCACHE_MODE.ini /usr/local/etc/php/conf.d/opcache.ini

ARG CONTENT_SECURITY_POLICY=
ENV CONTENT_SECURITY_POLICY=${CONTENT_SECURITY_POLICY}

# Configure Apache.

# first remove the default configuration
RUN rm /etc/apache2/sites-available/*.conf && \
    rm /etc/apache2/sites-enabled/*.conf

# Then add the WissKI site
ADD apache.d/conf/ports.conf /etc/apache2/ports.conf
ADD apache.d/sites-available/wisski.conf /etc/apache2/sites-available/wisski.conf

# And enable it
RUN a2ensite wisski

# volumes for composer
VOLUME /var/www/.composer
VOLUME /var/www/data

# Add and configure the entrypoint
ADD scripts/entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD ["apache2-foreground"]

# Add the user_shell.sh
ADD scripts/user_shell.sh /user_shell.sh
ADD ssh/ /ssh/
VOLUME /ssh/hostkeys/
RUN chmod 700 /ssh/keys.sh && \
    chmod 700 /ssh/start.sh && \
    chmod 777 /user_shell.sh && \
    chsh www-data --shell /user_shell.sh

# expose port 8080
EXPOSE 8080