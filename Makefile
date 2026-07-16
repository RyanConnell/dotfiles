install:
	(cd dotter && go run cmd/dotter.go install -c ../config.yaml -a ../apps -o ../build)

render:
	(cd dotter && go run cmd/dotter.go install -c ../config.yaml -a ../apps -o ../build --no-stow)

diff:
	(cd dotter && go run cmd/dotter.go diff -c ../config.yaml -a ../apps -o ../build)
