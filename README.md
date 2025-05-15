# blogging-example
Coding Challenge Solution

# How to run

Clone and CD into repo:

```bash
git clone git@github.com:youming-lin/blogging-example.git
cd blogging-example
```

Run main.go with DB connection:
```bash
export DB_URI=$(realpath data.sqlite) && go run main.go
```

# Additional work needed
1. Use a real database
2. Authentication and authorization
3. Dockerfile
4. Unit and integration tests
5. CI/CD pipelines
