node ./playground/update-relayer-env.js

export $(cat ./playground/relayers.env | xargs)
# echo $betchain_transfer_channel
# echo $gaia_transfer_channel
echo $betchain_oracle_channel
echo $bandchain_oracle_channel


# bccli tx coinpricebet \
# set-channel band-cosmoshub transfer $betchain_transfer_channel \
# --from validator --keyring-backend test -y -b block

# bccli tx coinpricebet \
# set-channel band-consumer transfer $gaia_transfer_channel \
# --from validator --keyring-backend test -y -b block

bccli tx coinpricebet \
set-channel ibc-bandchain coinpricebet $betchain_oracle_channel \
--from validator --keyring-backend test -y -b block

bccli tx coinpricebet \
set-channel band-consumer oracle $bandchain_oracle_channel \
--from validator --keyring-backend test -y -b block
