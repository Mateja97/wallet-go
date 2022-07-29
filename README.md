# wallet-go
Wallet for funds

# Run 

go run ./cmd/wallet/ -kafka.topic=wallet.public.wallet -kafka.brokers="172.17.0.1:9092" -db.host=172.17.0.1 -db.port=:5432 -db.usr=wallet -db.pw=wallet -db.name=wallet 