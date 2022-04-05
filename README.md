# URL-shortener-Dcard-backend-intern

## Tools

- Programming language: Go
- RDBMS: Sqlite
- Cache: Redis
- ORM lib: gorm
- Redis lib: go-redis

## Process

1. Generate the shortened URL

- Check request combines the required
- Check the "expireAt" field to make sure that is not expired. otherwise, return the bad request of status code
- Generate the UUID for the unique
- Use base64 to encrypt the UUID to be a shortened URL
- Store the shortened URL into the cache and database

2. Get original URL

- Check request is combine the required
- Query from the cache
   - If data existed
     - redirect the original URL
- Query from the database
   - If data existed
     - check the "expireAt" field to make sure that is not expired
       - If expired
         - Return Not Found of status code
       - Else
         - redirect the original URL
   - Else
     - Return Not Found of status code

## TODO

- Create the cronjob to delete the expired rows(Database)
- The verify work by application not the cache directly
