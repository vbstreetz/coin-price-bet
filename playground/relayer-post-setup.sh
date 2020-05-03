echo "typical abstract shoe junior annual idle conduct extend high source cliff zero quality brick fluid spare roast pulp claw swear bicycle lens teach digital" | bccli keys add vb --recover --keyring-backend test

bccli tx send requester $(bccli keys show -a vb) 1000000000000stake --keyring-backend test -y

# ufw allow 26657
