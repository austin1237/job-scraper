# AWS go lambdas running on provided.al2 runtime have to be called bootstrap 
# https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-naming

packageLambdas: packageScraper packageProxy packageJobNotifier
	
packageScraper:
	cd scraper && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls

packageProxy:
	cd proxy && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls

packageJobNotifier:
	cd jobNotifier && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls


