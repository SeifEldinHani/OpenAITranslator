## Running the Project

To start the project, use the following command:

```bash
docker-compose up --build
```

To Run the Project's tests, use the following command: 
```bash
docker-compose up -d && docker-compose exec app go test ./...
```
