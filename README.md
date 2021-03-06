# WKHTMLTOPDF as a service

Print anything to a PDF as a Golang webservice. The resulting image is Ubuntu-based, around 250MB in size and uses the latest version of wkhtmltopdf (0.12.6).

## Usage

Build the image:

`docker build . -t nwo/wkhtmltopdf`

Run the container. Mount a shared directory to `/app/shared` where the PDF will get stored:

`docker run --rm -e APP_HOST=:3000 -v $(pwd)/shared:/app/shared -p 3000:3000 nwo/wkhtmltopdf`

Call the service (i.e. with [httpie](https://httpie.org/)):

`http POST localhost:3000 options='-q -s A4 -B 0.5in -L 0.5in -R 0.5in -T 0.5in --encoding UTF-8 --title "My Document" --load-error-handling ignore' type=file file=input.html`

If you use `type: file` place the input file in the shared dir before calling the service.

## Payload

```js
{
    "options": "", // a string of wkhtmltopdf options
    "type": "", // can be one of 'file', 'string' or 'url'
    "file": "", // if the type is 'file', provide the input file here. it has to reside in /app/shared.
    "string": "", // if the type is string, provide the input string here.
    "url": "", // you now what to do.
}
```
