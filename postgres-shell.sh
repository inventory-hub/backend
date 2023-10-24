docker exec -it $(docker ps -aqf "name=^inventory-hub-backend-fullstack-postgres-1$") psql -h localhost -U calin -d inventory-hub-db -W

