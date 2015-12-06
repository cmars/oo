
PROJECT_ROOT := $(shell cd ../../..; pwd)

all: files/oostore build

build:
	JUJU_REPOSITORY=$(PROJECT_ROOT)/charms charm build

files/oostore: $(PROJECT_ROOT)/bin/oostore
	mkdir -p files
	cp $< $@

$(PROJECT_ROOT)/bin/oostore:
	cd $(PROJECT_ROOT); gb build

clean:
	$(RM) files/oostore
	$(RM) -r $(PROJECT_ROOT)/charms/trusty/oostore

.PHONY: all compose

