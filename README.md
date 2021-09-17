## Scrapper:
this is a small tool to scrap some data from a provided url

## Cloning

```
https://github.com/rakib32/go-scrapper.git
cd go-scrapper
```


## Building and Running

### Without Docker

```bash
cd go-scrapper
go build 
./go-scrapper -url  "http://www.google.com"
```

### With Docker

#### Local Development
```
cd go-scrapper
make develop
```

#### Production Build
* Build the image
```bash
cd go-scrapper
make build (It creates a image with go-scrtapper:latest tag)
```
* Run the image with custom args

```bash

docker run -it go-scrapper:latest /app -url "http://www.google.com"
or 
make run (Just change the args in makefile)
```

## usage

```
./go-scrapper [flags]
```

with the flags being
```
    -url="http://www.google.com": Url to scrap
    
```
for example
```
./go-scrapper -u "http://www.google.com"
```

## Testing
Use following command to run the test
```
go test ./...
OR
make test
```
