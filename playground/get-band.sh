export $(cat ./playground/relayers.env | xargs)

# Get some band for the bandchain relayer
# The relayer pays for the oracle requests
curl --location --request POST 'http://bandchain-ibc-hackathon.node.bandchain.org/faucet/request' \
--header 'Content-Type: application/json' \
--data-raw '{
	"address": "'"$(rly keys show ibc-bandchain)"'",
	"amount": 10000000
}'
