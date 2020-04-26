
# Create user-gaia user and request funds(atom) from cosmos faucet to it
make bccli o='keys add user-gaia --keyring-backend test'

user_gaia_addr=$(make bccli o='keys show -a user-gaia --keyring-backend test')
requester_addr=$(make bccli o='keys show -a requester --keyring-backend test')

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

#
export $(cat ./playground/relayers.env | xargs)

#

make bccli o="tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $requester_addr \
5000000000transfer/$betchain_transfer_channel/uatom \
--from user-gaia \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
--keyring-backend test \
-b block \
-y"

#

make bccli o="tx \
coinpricebet buy 1000000000transfer/$betchain_transfer_channel/uatom \
--from requester --keyring-backend test -y -b block"

make bccli o="query bank balances $requester_addr"
