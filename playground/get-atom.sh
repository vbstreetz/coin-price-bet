export $(cat ./playground/relayers.env | xargs)

curl --location --request \
POST 'http://gaia-ibc-hackathon.node.bandchain.org:8000' \
--header 'Content-Type: application/javascript' \
--data-raw '{
 "address": "'"$(bccli keys show -a vb)"'",
 "chain-id": "band-cosmoshub"
}'

sleep 5

bccli tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $(bccli keys show -a vb) \
5000000000transfer/$betchain_transfer_channel/uatom \
--from vb \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
-b block \
-y

sleep 5

bccli query bank balances $(bccli keys show -a vb) \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id=band-cosmoshub
