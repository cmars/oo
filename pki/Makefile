
TOOLS_BIN := ../tools/bin

CFSSL=$(TOOLS_BIN)/cfssl

all: ca localhost.cert.json api.oo.cmars.tech.cert.json

ca: ca.csr ca-key.pem ca.pem

ca.csr ca-key.pem ca.pem: ca.json
	$(CFSSL) genkey -initca $< | cfssljson -bare ca

%.cert.json: %.csr.json
	$(CFSSL) gencert -hostname $* $< > $@

destroy: destroy-ca destroy-cert

destroy-ca:
	$(RM) ca.csr ca-key.pem ca.pem

destroy-cert:
	$(RM) *.cert.json

.PHONY: ca all destroy destroy-ca destroy-cert

