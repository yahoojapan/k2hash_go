#
# k2hash_go
#
# Copyright 2018 Yahoo Japan Corporation.
#
# Go driver for k2hash that is a NoSQL Key Value Store(KVS) library.
# For k2hash, see https://github.com/yahoojapan/k2hash for the details.
#
# For the full copyright and license information, please view
# the license file that was distributed with this source code.
#
# AUTHOR:   Hirotaka Wakabayashi
# CREATE:   Fri, 14 Sep 2018
# REVISION:
#

build:
	# use _build as the GOPATH-base	
	@/bin/echo "Running k2hash_go build"
	GOPATH=$(PWD)/_build go install -v github.com/yahoojapan/k2hash_go/... 
	@echo "OK - built the following binaries:"
	ls -l _build/bin	

init:
	# 0. remove the build directory
	rm -fr _build
	# 1. get source.
	@echo "Running k2hash_go init (fetching source code)"	
	git clone https://github.com/yahoojapan/k2hash_go.git _build/src/github.com/yahoojapan/k2hash_go
	# 2. syntax check.
	BAD_GOFMT_FILES=$(find ./_build -name '*.go' | xargs gofmt -l)
	@echo ".go files that are not gofmt-compliant (empty if all are fine): [$(BAD_GOFMT_FILES)]"
	# 3. install libk2hash
	sh utils/libk2hash.sh

test:
	@echo "Running k2hash_go test"
	GOPATH=$(PWD)/_build go test -v github.com/yahoojapan/k2hash_go/tests
	GOPATH=$(PWD)/_build go test -v github.com/yahoojapan/k2hash_go/tests -coverprofile=c.out
	GOPATH=$(PWD)/_build go tool cover -html=c.out

publish:
	@echo "Running k2hash_go publish"
	# TODO: add scp of binaries to Artifactory (or RPM package creation and uploading)

# Local Variables:
# c-basic-offset: 4
# tab-width: 4
# indent-tabs-mode: t
# End:
# vim600: noexpandtab sw=4 ts=4 fdm=marker
# vim<600: noexpandtab sw=4 ts=4

