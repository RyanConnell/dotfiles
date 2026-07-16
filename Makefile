install:
	(cd dotter && go run cmd/dotter.go install -c ../config.yaml -a ../apps -o ../build)

diff:
	(cd dotter && go run cmd/dotter.go diff -c ../config.yaml -a ../apps -o ../build)
