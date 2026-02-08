# file-sync

docker build -t peterjbishop/file-sync:latest .

<!-- configured to expose the container to my local network -->
docker run -d \
-v /Users/m4pro/Sync:/app/data \
-e WATCH_DIR=/app/data \
-e HOST=0.0.0.0 \
-e PORT=8080 \
 -p 8080:8080 peterjbishop/file-sync:latest