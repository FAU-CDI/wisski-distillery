FROM php:7-apache-buster
WORKDIR /var/www

# install and enable the various required php extension
RUN apt-get update && apt-get install -y \
        curl \
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
        soap \
        mysqli \
        opcache \
        pdo_mysql \
        xml \
        xmlrpc \
    && \
    pecl install ssh2-1.2 && \
    docker-php-ext-enable \
        curl \
        gd \
        intl \
        mysqli \
        opcache \
        pdo_mysql \
        soap \
        ssh2 \
        mysqli \
        xml \
        xmlrpc \
    && \
    docker-php-source delete

# enable the apache rewrite mod
RUN a2enmod rewrite

# install composer and add it to path
RUN curl -sS https://getcomposer.org/installer | php && \
    mv composer.phar /usr/local/bin/composer
ENV PATH "/usr/local/bin:/var/www/data/project/vendor/bin:$PATH"

# remove default configuration
RUN rm /etc/apache2/sites-available/*.conf && \
    rm /etc/apache2/sites-enabled/*.conf

ADD patch/easyrdf.patch /patch/easyrdf.patch

# Add wisski configuration
ADD conf/ports.conf /etc/apache2/ports.conf
ADD conf/wisski.conf /etc/apache2/sites-available/wisski.conf
ADD conf/wisski.ini /usr/local/etc/php/conf.d/wisski.ini
RUN a2ensite wisski

# volumes for composer
VOLUME /var/www/.composer
VOLUME /var/www/data

# Add and configure the entrypoint
ADD scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["apache2-foreground"]

# Add the provision script
ADD scripts/provision_container.sh /provision_container.sh

# Add the user_shell.sh
ADD scripts/user_shell.sh /user_shell.sh

# expose port 8080
EXPOSE 8080