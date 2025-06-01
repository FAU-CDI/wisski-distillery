if [[ "$USER" -ne "www-data" ]]; then
    return
fi

export "PATH=/var/www/data/project/vendor/bin:$PATH"
