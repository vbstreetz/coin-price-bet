echo "typical abstract shoe junior annual idle conduct extend high source cliff zero quality brick fluid spare roast pulp claw swear bicycle lens teach digital" | bccli keys add vb --recover --keyring-backend test

# Get some stake
bccli tx send requester $(bccli keys show -a vb) 1000000000000stake --keyring-backend test -y
# bccli tx send requester cosmos1crecwthf7rvddf3vznsrw4wwlq4pspkp0ehgkv 1000000stake --keyring-backend test -y --home /root/.bccli

# ufw allow 26657

# Get some band for the bandchain relayer
# The relayer pays for the oracle requests
curl --location --request POST 'http://bandchain-ibc-hackathon.node.bandchain.org/faucet/request' \
--header 'Content-Type: application/json' \
--data-raw '{
	"address": "'"$(rly keys show ibc-bandchain)"'",
	"amount": 10000000
}'
