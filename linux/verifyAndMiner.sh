#!/bin/bash
OS="LINUX"
LinuxConfigDir="${HOME}/.usechain"
OSXConfigDir=${HOME}/Library/usechain
WindowsConfigDir=${HOME}/AppData/Roaming/usechain
ConfigDir="$LinuxConfigDir"

case "$OSTYPE" in
  solaris*) OS="SOLARIS" ;;
  linux*)   OS="LINUX" ;;
  bsd*)     OS="BSD" ;;
  darwin*)  OS="OSX" ; ConfigDir=$OSXConfigDir;; # Mac
  msys*)    OS="WINDOWS" ; ConfigDir=$WindowsConfigDir;;# Git Bash/msysGit
  cygwin*)  OS="WINDOWS" ; ConfigDir=$WindowsConfigDir;; # Cygwin
  *)        OS="UNKNOWN"; echo "unknown: $OSTYPE" ;;
esac

echo
echo "OS:$OS"
echo "Usechain ConfigDir: $ConfigDir"
if [ ${OS}x == "UNKNOWN"x ]; then
    echo "error OS $OSTYPE not supported, exit ..."
    exit;
fi

cur_dir=$(cd "$(dirname "$0")"; pwd)

# make config directory for user
if [ ! -d $ConfigDir ]; then
    echo "mkdir $ConfigDir"
    mkdir -p $ConfigDir
fi


autoMiner="./autoMiner"
if [ ! -f "$autoMiner" ];then
	go build ./autoMiner.go
fi

count=`ps -ef|grep "./used"|grep -v grep|wc -l`
if [ $count -eq 0 ]; then
        echo "start used"
        nohup ./used --rpc --rpcaddr="0.0.0.0" --rpccorsdomain "*"  --rpcapi="web3,eth,admin,personal,miner,net" --gasprice "10"  --txpool.pricelimit "10"  --verbosity 3  &
fi

userCert="$ConfigDir/user.crt"
if [ ! -f "$userCert" ]; then

userDataFile="$ConfigDir/userData.json"

#check userdata.json
if [ ! -f "$userDataFile" ]; then
    ./used verify
else
	#if userdata.json exist, do you need modify
	read -r -p ""  -t 0.01 inputx
	read -r -p "Your userData.json already exit, do you need to modify?[Y/n] " input
	case $input in
		    [yY][eE][sS]|[yY])
				echo "Yes"
				./used verify
				;;

		    [nN][oO]|[nN])
				echo "No"
	       			;;
			    *)

			echo "Invalid input..."
			;;
	esac
fi
#get idkey
idkey=`./used verify --info=userData.json --photo="use.jpg;use2.jpg"`
echo $idkey
function get_idKey()
{
  local json=$1
  local key=$2

  if [[ -z "$3" ]]; then
    local num=1
  else
    local num=$3
  fi

  local value=$(echo "${json}" | awk -F"[,:}]" '{for(i=1;i<=NF;i++){if($i~/'${key}'\042/){print $(i+1)}}}' | tr -d '"' | sed -n ${num}p)

  echo ${value}
}

idk=$(get_idKey "$idkey" 'idKey')

if [ -z "$idk" ]; then
	echo "Cannot get idkey"
	exit 1
fi

#get certification
cert=`./used verify --query=$idk`
if [ -z "$cert" ]; then
        echo "Your information is verifying, Cannot get cert for now"
	exit 1
fi
fi
echo "You already have user.crt, you can register and miner directly, please continue..."

#autominer include register miner, start miner
sleep 3
./autoMiner