#!/bin/bash
autoMiner=".\autoMiner"
if [ ! -f "$autoMiner" ];then
	go build ./autoMiner.go
fi

start usedStart.bat

userCert="$HOME\AppData\Roaming\usechain\user.crt"
if [ ! -f "$userCert" ]; then

#check userdata.json
userDataFile="$HOME\AppData\Roaming\usechain\userData.json"
if [ ! -f "$userDataFile" ]; then
    ./used.exe verify
else
	#if userdata.json exist, do you need modify
	read -r -p ""  -t 0.01 input
	read -r -p "Your userData.json already exit, do you need to modify?[Y/n] " input
	case $input in
		    [yY][eE][sS]|[yY])
				echo "Yes"
				used.exe verify
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
idkey=`used.exe verify --info=userData.json --photo="use.jpg;use2.jpg"`
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
#cert=`used.exe verify --query=$idk`
#if [ -z "$cert" ]; then
#        echo "Your information is verifying, Cannot get cert for now"
#	exit 1
#fi
fi

echo "You already have user.crt, you can register identity and miner directly, continue..."
#autoMiner include register miner, start miner
./autoMiner.exe
sleep 15
