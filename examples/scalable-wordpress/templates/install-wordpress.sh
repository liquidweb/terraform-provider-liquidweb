groupadd ${user}
useradd -g ${user} ${user}

mkdir -p /home/${user}/www
chmod g+x /home/${user}
usermod nginx -G wordpress

curl --silent https://wordpress.org/latest.tar.gz|tar -C /home/${user}/www --strip-components=1 -xz 'wordpress/*'

mv /root/wp-config.php /home/${user}/www/

chown -R wordpress. /home/wordpress/www