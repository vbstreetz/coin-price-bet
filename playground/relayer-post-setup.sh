
# Create user-gaia user and request funds(atom) from cosmos faucet to it
# Then send some to vb
make bccli o='keys add user-gaia --keyring-backend test'

user_gaia_addr=$(make bccli o='keys show -a user-gaia --keyring-backend test')
vb_addr=$(make bccli o='keys show -a vb --keyring-backend test')

curl --location --request \
POST 'http://gaia-ibc-hackathon.node.bandchain.org:8000' \
--header 'Content-Type: application/javascript' \
--data-raw '{
 "address": "'"$user_gaia_addr"'",
 "chain-id": "band-cosmoshub"
}'

make bccli o="query bank balances $user_gaia_addr \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id=band-cosmoshub"

make bccli o="tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $vb_addr \
5000000000transfer/$betchain_transfer_channel/uatom \
--from user-gaia \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
--keyring-backend test \
-b block \
-y"

make bccli o="tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $vb_addr \
5000000000xxx \
--from user-gaia \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
--keyring-backend test \
-b block \
-y"

#
export $(cat ./playground/relayers.env | xargs)

# Request band via bandchain lite blockchain relayer
# Then send some to vb

curl --location --request POST 'http://bandchain-ibc-hackathon.node.bandchain.org/faucet/request' \
--header 'Content-Type: application/json' \
--data-raw '{
	"address": "'"$user_gaia_addr"'",
	"amount": 10000000
}'

make bccli o="tx oracle transfer \
transfer $band_transfer_channel \
10000000 $vb_addr \
5000000000transfer/$betchain_transfer_channel/uband \
--from user_gaia_addr \
--node http://bandchain-ibc-hackathon.node.bandchain.org:26657 \
--chain-id ibc-bandchain \
--keyring-backend test \
-b block \
-y"

#

make bccli o="tx \
coinpricebet buy 1000000000transfer/$betchain_transfer_channel/uatom \
--from requester --keyring-backend test -y -b block"

make bccli o="query bank balances $vb_addr"
make bccli o="query bank balances $user_gaia_addr"

#

make bccli o="query coinpricebet latest-coin-prices 0"

# ufw allow 26657
