#!/bin/sh
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

# Sets the default locale. LC_ALL has precedence over other LC* variables.
unset LANG
unset LANGUAGE
LC_ALL=en_US.utf8
export LC_ALL

# Sets PATH. setup_*.sh uses useradd command
PATH=${PATH}:/usr/sbin:/sbin

# an unset parameter expansion will fail
set -u

# umask 022 is enough
umask 022

# environments
SRCDIR=$(cd $(dirname "$0") && pwd)
DEBUG=1
if test "${DEBUG}" -eq 1; then
	TAG="$(basename $0) -s"
else
	TAG=$(basename $0)
fi
USER=$(whoami)
LOGLEVEL=info

# Checks if k2hash is installed 
#
# Params::
#   no params
#
# Returns::
#   0 on installed
#   1 on not installed
#
which_k2hash() {
	which k2hlinetool >/dev/null 2>&1
	if test "${?}" = "0"; then
		logger -t ${TAG} -p user.info "k2hash already installed"
		return 0
	fi
	return 1
}

# Determines the current OS
#
# Params::
#   no params
#
# Returns::
#   0 on success
#   1 on failure
#
setup_os_env() {
	if test -f "/etc/os-release"; then
		. /etc/os-release
		OS_NAME=$ID
		OS_VERSION=$VERSION_ID
	else
		logger -t ${TAG} -p user.warn "unknown OS, no /etc/os-release and /etc/centos-release falling back to CentOS-7"
		OS_NAME=centos
		OS_VERSION=7
	fi

	if test "${OS_NAME}" = "ubuntu"; then
		logger -t ${TAG} -p user.notice "ubuntu configurations are currently equal to debian one"
		OS_NAME=debian
	elif test "${OS_NAME}" = "centos"; then
		if test "${OS_VERSION}" != "7"; then
			logger -t ${TAG} -p user.err "centos7 only currently supported, not ${OS_NAME} ${OS_VERSION}"
			exit 1
		fi
	fi

	HOSTNAME=$(hostname)
	logger -t ${TAG} -p user.debug "HOSTNAME=${HOSTNAME} OS_NAME=${OS_NAME} OS_VERSION=${OS_VERSION}"
}

# Builds k2hash from source code
#
# Params::
#   $1 os_name
#
# Returns::
#   0 on success
#   1 on failure(exit)
#
make_k2hash() {

	_os_name=${1:?"os_name should be nonzero"}

	if test "${_os_name}" = "debian" -o "${_os_name}" = "ubuntu"; then
		_configure_opt="--with-gcrypt"
		sudo apt-get update -y
		sudo apt-get install -y git curl autoconf autotools-dev gcc g++ make gdb libtool pkg-config libyaml-dev libgcrypt20-dev
	elif test "${_os_name}" = "fedora"; then
		_configure_opt="--with-nss"
		sudo dnf install -y git curl autoconf automake gcc gcc-c++ gdb make libtool pkgconfig libyaml-devel nss-devel
	elif test "${_os_name}" = "centos" -o "${_os_name}" = "rhel"; then
		_configure_opt="--with-nss"
		sudo yum install -y git curl autoconf automake gcc gcc-c++ gdb make libtool pkgconfig libyaml-devel nss-devel
	else
		logger -t ${TAG} -p user.error "OS should be debian, ubuntu, fedora, centos or rhel"
		exit 1
	fi

	logger -t ${TAG} -p user.debug "git clone https://github.com/yahoojapan/k2hash"
	git clone https://github.com/yahoojapan/k2hash
	cd k2hash
	
	logger -t ${TAG} -p user.debug "git clone https://github.com/yahoojapan/fullock"
	git clone https://github.com/yahoojapan/fullock
	logger -t ${TAG} -p user.debug "git clone https://github.com/yahoojapan/k2hash"
	git clone https://github.com/yahoojapan/k2hash
	
	if ! test -d "fullock"; then
		echo "no fullock"
		exit 1
	fi
	cd fullock
	./autogen.sh
	./configure --prefix=/usr
	make
	sudo make install
	
	if ! test -d "../k2hash"; then
		echo "no k2hash"
		exit 1
	fi
	cd ../k2hash
	./autogen.sh
	./configure --prefix=/usr ${_configure_opt}
	make
	sudo make install
	
	return 0
}

#
# main loop
#

setup_os_env

which_k2hash
if test "${?}" = "0"; then
	exit 0
fi

if test "${OS_NAME}" = "fedora"; then
	which sudo
	if test "${?}" = "1"; then
		dnf install -y sudo
	fi
	if test "${?}" = "1"; then
		sudo dnf install -y bc
	fi
	which bc
	if test "${?}" = "1"; then
		sudo dnf install -y bc
	fi
	if test "${OS_VERSION}" = "28"; then
		sudo dnf install -y curl
		curl -s https://packagecloud.io/install/repositories/antpickax/stable/script.rpm.sh | sudo bash
		sudo dnf install  k2hash-devel -y
	elif test "${OS_VERSION}" = "29"; then
		sudo dnf install -y curl
		curl -s https://packagecloud.io/install/repositories/antpickax/current/script.rpm.sh | sudo bash
		sudo dnf install  k2hash-devel -y
	else
		make_k2hash ${OS_NAME}
	fi
elif test "${OS_NAME}" = "debian" -o "${OS_NAME}" = "ubuntu"; then
	which sudo
	if test "${?}" = "1"; then
		apt-get update -y
		apt-get install -y sudo
	fi
	which bc
	if test "${?}" = "1"; then
		sudo apt-get install -y bc
	fi
	if test "${OS_VERSION}" = "9"; then
		sudo apt-get install -y curl
		curl -s https://packagecloud.io/install/repositories/antpickax/stable/script.deb.sh | sudo bash
		sudo apt-get install -y k2hash-dev
	else
		make_k2hash ${OS_NAME}
	fi
elif test "${OS_NAME}" = "centos" -o "${OS_NAME}" = "rhel"; then
	which sudo
	if test "${?}" = "1"; then
		yum install -y sudo
	fi
	which bc
	if test "${?}" = "1"; then
		sudo yum install -y bc
	fi
	yum install -y k2hash-devel
	if test "${?}" != "0"; then
		sudo yum install -y curl
		curl -s https://packagecloud.io/install/repositories/antpickax/stable/script.rpm.sh | bash
		sudo yum install -y k2hash-devel
		if test "${?}" != "0"; then
			make_k2hash ${OS_NAME}
		fi
	fi
else
	echo "[NO] OS must be either fedora or centos or debian or ubuntu, not ${OS_NAME}"
	exit 1
fi

if ! test -f "/dev/log"; then
	if test -f "/run/systemd/journal/dev-log"; then
		sudo ln -s /run/systemd/journal/dev-log /dev/log
	else
		echo "/dev/log unavailable"
	fi
fi

which_k2hash
if test "${?}" = "0"; then
	exit 0
fi

exit 0

# Local Variables:
# c-basic-offset: 4
# tab-width: 4
# indent-tabs-mode: t
# End:
# vim600: noexpandtab sw=4 ts=4 fdm=marker
# vim<600: noexpandtab sw=4 ts=4
