build:
	docker build . -t nwo/wkhtmltopdf

run:
	docker run --rm -e APP_HOST=:3000 -v $(shell pwd)/shared:/app/shared -p 3000:3000 nwo/wkhtmltopdf
