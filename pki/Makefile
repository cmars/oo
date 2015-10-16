
TOOLS_BIN := ../tools/bin

CFSSL=$(TOOLS_BIN)/cfssl
CFSSLJSON=$(TOOLS_BIN)/cfssljson
JQ=jq

all: ca \
	localhost.pem localhost-key.pem \
	api.oo.cmars.tech.pem api.oo.cmars.tech-key.pem

ca: ca.csr ca-key.pem ca.pem

ca.csr ca-key.pem ca.pem: ca.json
	$(CFSSL) genkey -initca $< | $(CFSSLJSON) -bare ca

%.cert.json: %.csr.json
	$(CFSSL) gencert -hostname $* $< > $@

%.pem %-key.pem: %.cert.json
	$(JQ) -r .cert $*.cert.json > $*.pem
	cat ca.pem >> $*.pem
	$(JQ) -r .key $*.cert.json > $*-key.pem

destroy: destroy-ca destroy-cert

destroy-ca:
	$(RM) ca.csr ca-key.pem ca.pem

destroy-cert:
	$(RM) *.cert.json localhost*.pem api.oo*.pem

.PHONY: ca all destroy destroy-ca destroy-cert

