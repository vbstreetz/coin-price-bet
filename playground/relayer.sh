# Remove old config
rm -rf ~/.relayer

rly config init

# Add config after these commands can check your config file at `~/.relayer/config/config.yaml`
rly chains add -f ./relayer/gaia.json
rly chains add -f ./relayer/ibc-bandchain.json
rly chains add -f ./relayer/betchain.json

# Add relayer account (Recover by mnemonic help for developing)
rly keys restore band-consumer relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"
rly keys restore ibc-bandchain relayer "mix swift essence lawsuit plastic major social copper chicken aisle caution unfold leaf turtle prize remove gravity tourist gym parade number street twelve long"
rly keys restore band-cosmoshub relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"

# Update default relayer for each chain
rly ch edit band-consumer key relayer
rly ch edit ibc-bandchain key relayer
rly ch edit band-cosmoshub key relayer

# And make sure every relayer have default coin in each chain by
rly q bal band-consumer
rly q bal ibc-bandchain
rly q bal band-cosmoshub

# Init lite client and save state for each chain
rly lite init band-consumer -f
rly lite init ibc-bandchain -f
rly lite init band-cosmoshub -f

# Create path(specific connection between chain)
rly pth gen  band-consumer transfer band-cosmoshub transfer transfer
rly pth gen  band-consumer coinpricebet ibc-bandchain oracle oracle

# Create connection and channel from path
echo "" > ./playground/relayer-create.log
rly tx link transfer >> ./playground/relayer-create.log
rly tx link oracle >> ./playground/relayer-create.log

./playground/update-env.sh

## Seperate run these command in different windows
#rly st transfer --debug
#rly st oracle --debug


