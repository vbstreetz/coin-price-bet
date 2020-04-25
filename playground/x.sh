export $(cat ./playground/relayers.env | xargs)

make bccli o="tx coinpricebet request-gold-price-update --from requester -b block -y"
