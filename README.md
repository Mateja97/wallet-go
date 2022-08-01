# wallet-go
Wallet for funds

# Preperation 
 - docker-compose up -d
 - docker-compose exec walletdb sh -c 'psql -U wallet < /db/script.sql'
 - debezium connector through POST on http://<host ip>/connectors/ 
with body: 
{
    "name": "players-connector",
    "config": {
        "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
        "database.hostname":-host ip-,
        "database.port": "5432",
        "database.user": "wallet",
        "database.password": "wallet",
        "database.dbname": "wallet",
        "database.server.name": "wallet",
        "plugin.name": "pgoutput",
        "table.include.list": "public.wallet"
    }
}

# Run 
go run ./cmd/wallet/ -kafka.topic=wallet.public.wallet -kafka.brokers="-host ip-:9092" -db.host=-host ip- -db.port=:5432 -db.usr=wallet -db.pw=wallet -db.name=wallet 