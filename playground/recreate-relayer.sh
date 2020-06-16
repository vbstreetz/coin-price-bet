set -e

rly lite delete band-consumer
rly lite init band-consumer -f

#rly pth delete transfer
rly pth delete oracle

#rly pth gen band-consumer transfer band-cosmoshub transfer transfer
rly pth gen band-consumer coinpricebet ibc-bandchain oracle oracle

echo "linking paths..."
echo "" > ./playground/relayer-create.log
#rly tx link transfer >> ./playground/relayer-create.log
rly tx link oracle >> ./playground/relayer-create.log

./playground/update-env.sh
