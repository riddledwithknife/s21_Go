all:
	cd cmd/candy-server-server && go run main.go --tls-certificate=../../ca/localhost/cert.pem --tls-key=../../ca/localhost/key.pem --tls-ca=../../ca/minica.pem --scheme=https --tls-port=3333

test:
	curl --key ca/client/key.pem --cert ca/client/cert.pem --cacert ca/minica.pem -XPOST -H "Content-Type: application/json" -d '{"money": 20, "candyType": "AA", "candyCount": 1}' "https://localhost:3333/buy_candy"