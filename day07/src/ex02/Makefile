godoc_install:
	go install golang.org/x/tools/cmd/godoc@latest

godoc_start:
	~/go/bin/godoc -http=:6060

godoc_html:
	../../../../go/bin/godoc -url "http://localhost:6060/pkg/day07/ex00" > doc.html

zip_html:
	zip docs.zip doc.html
	rm doc.html

unzip_html:
	unzip docs.zip
	rm docs.zip

clean:
	rm doc.html docs.zip