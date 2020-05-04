export $(cat ./playground/relayers.env | xargs)

addr=cosmos13dus5s0c2ey32509xy6hgk8s923wufrv2dyy6g
chain_id=band-consumer

# Get the sequence and account numbers for `requester` to construct the below requests
curl -s http://localhost:1318/auth/accounts/$addr

account_number=9
sequence=0

# Create buy_gold raw transaction
curl -XPOST -s http://localhost:1318/coinpricebet/buy --data-binary '{"base_req":{"from":"'$addr'","chain_id":"'$chain_id'"},"amount":"1000000000transfer/'$betchain_transfer_channel'/uatom"}' > playground/unsigned-tx.json

# Then sign the transaction
bccli tx sign playground/unsigned-tx.json --from requester --offline --chain-id $chain_id --sequence $sequence --account-number $account_number > playground/signed-tx.json

bccli tx \
coinpricebet buy 1000000000transfer/$betchain_transfer_channel/uatom \
--from vb \
-b block \
-y

# And finally broadcast the signed transaction
bccli tx broadcast playground/signed-tx.json

# Confirm
bccli query bank balances $addr
