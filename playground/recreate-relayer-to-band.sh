rly lite delete band-consumer
rly lite init band-consumer -f
rly pth delete oracle
rly pth gen band-consumer coinpricebet ibc-bandchain oracle oracle
rly tx link oracle > ./playground/relayer-create.log


node ./playground/update-relayer-env.js


export $(cat ./playground/relayers.env | xargs
make bccli o="tx coinpricebet \
set-channel ibc-bandchain coinpricebet $betchain_oracle_channel \
--from validator --keyring-backend test -y -b block"
make bccli o="tx coinpricebet \
set-channel band-cosmoshub transfer $betchain_transfer_channel \
--from validator --keyring-backend test -y -b block"
