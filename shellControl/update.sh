#!/bin/bash
p="/root/server/"
#文件分发到版本控制
function GitFile(){
  path=$p$1"/"
  if [ ! -d "$path" ]; then
    mkdir -p $p
    cd $p
    git clone git@gitee.com:ballwang/$1.git
    cd $path
    git remote set-url origin git@gitee.com:ballwang/$1.git
    git config --global user.name "Ballwang"
    git config --global user.email "Ballwang@foxmail.com"
    git pull
    git add .
    git commit -m "服务修改"
    git push
  else
    mv -f /root/go/src/github.com/Ballwang/tugo/$1 /root/server/$1/
    \cp -rf /root/go/src/github.com/Ballwang/tugo/config /root/server/$1/
    cd $path
    git pull
    git add .
    git commit -m "服务修改"
    git push
  fi
}

if [ ! -n "$1" ] ;then
    cd /root/go/src/github.com/Ballwang/tugo/
    git pull
    for file in ` ls /root/go/src/github.com/Ballwang/tugo/ |grep .go `
    do
        echo $file
    cd /root/go/src/github.com/Ballwang/tugo/
    go build $file
    fileName=${file%%.go*}
    GitFile $fileName
    done
else
    x=$1
    name=${x%%:*}
    OLD_IFS="$IFS"
    IFS=","
    array=($name)
    IFS="$OLD_IFS"
        cd /root/go/src/github.com/Ballwang/tugo/
    git pull
    for eachName in ${array[*]}
    do
          cd /root/go/src/github.com/Ballwang/tugo/
          echo $eachName".go"
          go build $eachName".go"
      GitFile $eachName
    done
fi