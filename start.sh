read -p "are you sure? " -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
  rm -rf ~/.bc*

  make bcd o='init validator --chain-id band-consumer'
  make bccli o='keys add validator --keyring-backend test'
  echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" | make bccli o='keys add requester --recover --keyring-backend test'
  echo "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice" | make bccli o='keys add relayer --recover --keyring-backend test'

  make bcd o='add-genesis-account validator 10000000000000stake --keyring-backend test'
  make bcd o='add-genesis-account requester 10000000000000stake --keyring-backend test'
  make bcd o='add-genesis-account relayer 10000000000000stake --keyring-backend test'

  make bccli o='config chain-id band-consumer'
  make bccli o='config output json'
  make bccli o='config indent true'
  make bccli o='config trust-node true'
  make bccli o='config keyring-backend test'

  make bcd o='gentx --name validator --keyring-backend test'
  make bcd o='collect-gentxs'

  # Run chain
  make bcd o='start --rpc.laddr=tcp://0.0.0.0:26657 --pruning=nothing'
fi