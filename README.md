Project: URL Shortener

Instructions: 
    - To get the project to work, you will need to change your sql connection string in internal/private.go. 

Notes to self: 
    - Check test code coverage: 
        go test -cover ./...
    - See Testing in HTML: 
        go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
    - using sqlmock 
