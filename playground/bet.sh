export $(cat ./playground/relayers.env | xargs)

m=m2

echo "shaft canyon strategy raccoon silly learn pole gospel fee turtle right opinion wise peace cabin melody anger library rice soul undo broom pig left" | bccli keys add $m --recover --keyring-backend test

addr=$(bccli keys show -a $m --keyring-backend test)

bccli tx send requester $addr 100000000stake --keyring-backend test -y

curl https://witnet.tools/vb-rest/auth/accounts/$addr

chain_id=band-consumer
account_number=0
sequence=0

# Create buy_gold raw transaction
curl -XPOST -s https://witnet.tools/vb-rest/coinpricebet/place-bet --data-binary '{"base_req":{"from":"'$addr'","chain_id":"'$chain_id'"},"amount":"1000000transfer/'$betchain_transfer_channel'/uatom"}' > playground/unsigned-tx.json

# Then sign the transaction
bccli tx sign playground/unsigned-tx.json --from $m --offline --chain-id $chain_id --sequence $sequence --account-number $account_number > playground/signed-tx.json

# get atom

curl --location --request \
POST 'http://gaia-ibc-hackathon.node.bandchain.org:8000' \
--header 'Content-Type: application/javascript' \
--data-raw '{
 "address": "'$addr'",
 "chain-id": "band-cosmoshub"
}'

bccli tx transfer transfer \
transfer $gaia_transfer_channel \
10000000 $addr \
5000000000transfer/$betchain_transfer_channel/uatom \
--from $m \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id band-cosmoshub \
-b block \
-y

bccli query bank balances $(bccli keys show -a vb) \
--node http://gaia-ibc-hackathon.node.bandchain.org:26657 \
--chain-id=band-cosmoshub

# place bet

bccli tx \
coinpricebet place-bet 1000000000transfer/$betchain_transfer_channel/uatom btc  \
--from $m \
-b block \
-y

# And finally broadcast the signed transaction
bccli tx broadcast playground/signed-tx.json
