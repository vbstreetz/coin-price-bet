cp -f ./nginx/nginx.proxy.conf /usr/local/etc/nginx/servers/vbstreetz-chain.conf
sudo launchctl unload /Library/LaunchAgents/homebrew.mxcl.nginx.plist || xargs echo
sudo launchctl load /Library/LaunchAgents/homebrew.mxcl.nginx.plist || xargs echo
echo 'done'
