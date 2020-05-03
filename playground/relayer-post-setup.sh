vb_addr=$(make bccli o='keys show -a vb --keyring-backend test')
vb_band_addr=$(rly keys show vbstreetz)

bandchain_relayer_addr=$(make bccli o='keys show -a band --keyring-backend test')
bandchain_relayer_band_addr=band16q6l6yp05mwtvxlcd03a9hlj8tka97ah7kvpuu # $(rly keys show ibc-bandchain)

export $(cat ./playground/relayers.env | xargs)

# Get some atom for vb

curl --location --request \
POST 'http://gaia-ibc-hackathon.node.bandchain.org:8000/' \
--header 'Content-Type: application/javascript' \
--data-raw '{
 "address": "'"$vb_addr"'",
 "chain-id": "band-cosmoshub"
}'

make bccli o="query bank balances $vb_addr \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id=band-cosmoshub"

make bccli o="tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $vb_addr \
5000000000transfer/$betchain_transfer_channel/uatom \
--from vb \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
--keyring-backend test \
-b block \
-y"

# Get some band for the bandchain relayer
# The relayer pays for the oracle requests

curl --location --request POST 'http://bandchain-ibc-hackathon.node.bandchain.org/faucet/request' \
--header 'Content-Type: application/json' \
--data-raw '{
	"address": "'"$bandchain_relayer_band_addr"'",
	"amount": 10000000
}'

#make bccli o="query bank balances $bandchain_relayer_addr \
#--node http://bandchain-ibc-hackathon.node.bandchain.org:26657 \
#--chain-id=ibc-bandchain"

#make bccli o="tx oracle transfer \
#transfer $band_transfer_channel \
#10000000 $vb_addr \
#5000000000transfer/$betchain_transfer_channel/uband \
#--from vb \
#--node http://bandchain-ibc-hackathon.node.bandchain.org:26657 \
#--chain-id ibc-bandchain \
#--keyring-backend test \
#-b block \
#-y"

#

make bccli o="tx \
coinpricebet buy 1000000000transfer/$betchain_transfer_channel/uatom \
--from requester --keyring-backend test -y -b block"

make bccli o="query bank balances $vb_addr"

make bccli o="query coinpricebet latest-coin-prices 0"

# ufw allow 26657
