# file-sync

docker build -t peterjbishop/file-sync:latest .
docker run \
-e HOST=0.0.0.0 \
-e PORT=8080 \
-p 8080:8080 peterjbishop/file-sync:latest