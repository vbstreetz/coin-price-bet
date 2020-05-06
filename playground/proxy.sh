ln -sf $PWD/nginx/nginx.proxy.conf ~/homebrew/etc/nginx/servers/betchain.conf
mkdir -p ~/homebrew/etc/supervisor.d
ln -sf $PWD/supervisor/supervisord.ini ~/homebrew/etc/supervisord.ini
ln -sf $PWD/supervisor/supervisor.d/dev/* ~/homebrew/etc/supervisor.d/
