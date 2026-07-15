install:
	(cd dotter && go run cmd/dotter.go install -c ../config.yaml -a ../apps -o ../build)
