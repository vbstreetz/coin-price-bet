echo "" > ./playground/relayer-create.log

rly lite delete vbstreetz
rly lite init vbstreetz -f

rly pth delete transfer
rly pth delete oracle

rly pth gen vbstreetz transfer band-cosmoshub transfer transfer
rly pth gen vbstreetz coinpricebet ibc-bandchain oracle oracle

rly tx link transfer >> ./playground/relayer-create.log
rly tx link oracle >> ./playground/relayer-create.log

./playground/update-env.sh

