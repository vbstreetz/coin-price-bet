export $(cat ./playground/relayers.env | xargs)

requester_addr=$(make bccli o='keys show -a requester --keyring-backend test')
make bccli o="query bank balances $requester_addr"


