docker run -d \
  --name searxng \
  -p 8888:8080 \
  -v "${PWD}/searxng:/etc/searxng" \
  -e "BASE_URL=http://localhost:8080/" \
  -e "INSTANCE_NAME=my-searxng" \
  --restart unless-stopped \
  searxng/searxng