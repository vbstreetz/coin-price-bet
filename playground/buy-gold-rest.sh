export $(cat ./playground/relayers.env | xargs)

requester_addr=$(bccli keys show -a requester --keyring-backend test)
chain_id=band-consumer

# Get the sequence and account numbers for `requester` to construct the below requests
curl -s http://localhost:1317/auth/accounts/$requester_addr

account_number=4
sequence=1

# Create buy_gold raw transaction
curl -XPOST -s http://localhost:1318/coinpricebet/buy --data-binary '{"base_req":{"from":"'$requester_addr'","chain_id":"'$chain_id'"},"amount":"1000000000transfer/'$betchain_transfer_channel'/uatom"}' > playground/unsigned-tx.json

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
bccli query bank balances $requester_addr

