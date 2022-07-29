docker-compose up -d
docker-compose exec walletdb sh -c 'psql -U wallet < /db/script.sql'
