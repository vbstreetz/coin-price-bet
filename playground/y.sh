export $(cat ./playground/relayers.env | xargs)

make bccli o="tx \
	coinpricebet buy 1000000000transfer/$betchain_transfer_channel/uatom \
	--from requester --keyring-backend test -y -b block"
