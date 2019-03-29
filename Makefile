# golang
GO      = CGO_ENABLE=0 go
GO-FMT  = ${GO} fmt
GO-VET  = ${GO} vet
GO-TEST = ${GO} test
GO-LIST = ${GO} list

# source code package list
PKGS = $(shell ${GO-LIST} ./...)

# target name label
TARGET-NAME = " ---> [$@]"

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    fmt      	run go fmt'
	@echo '    vet      	run go vet'
	@echo '    tests    	run go tests'
	@echo ''

fmt:
	@echo ${TARGET-NAME}
	@${GO-FMT} ${PKGS}

vet: fmt
	@echo ${TARGET-NAME}
	@${GO-VET} ${PKGS}

tests: vet
	@echo ${TARGET-NAME}
	@${GO-TEST} -count 10 ${PKGS}
