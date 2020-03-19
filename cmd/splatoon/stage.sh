#!/bin/bash
cd "$(dirname "$0")"
r="gachi"
t="now"
for OPT in "$@"
do
    case $OPT in
        -r)
            r="regular"
            ;;
        -l)
            r="league"
            ;;
        -n)
            t="next"
            ;;
    esac
#    shift
done
#echo $r $t $@
d=$(date "+%d" --utc)
h=$(date "+%H" --utc)
p=$((h%2+h/2))
if [ -e $d/$p$r$t ]
then
j=$(cat $d/$p$r$t)
else
if [ ! -e $d ]
then
rm $((d-1)) -r 2>/dev/null
mkdir -p $d
fi
echo Downloading...
j=$(curl -sH 'User-Agent: DiscordBot/1.0 (email public.yk@outlook.com)' https://spla2.yuu26.com/$r/$t)
echo $j > $d/$p$r$t
fi
#j='{"result":[{"rule":"ナワバリバトル","rule_ex":{"key":"turf_war","name":"ナワバリバトル","statink":"nawabari"},"maps":["ショッツル鉱山","タチウオパーキング"],"maps_ex":[{"id":17,"name":"ショッツル鉱山","image":"https://app.splatoon2.nintendo.net/images/stage/828e49a8414a4bbc0a5da3e61454ab148a9f4063.png","statink":"shottsuru"},{"id":8,"name":"タチウオパーキング","image":"https://app.splatoon2.nintendo.net/images/stage/96fd8c0492331a30e60a217c94fd1d4c73a966cc.png","statink":"tachiuo"}],"start":"2020-03-18T09:00:00","start_utc":"2020-03-18T00:00:00+00:00","start_t":1584489600,"end":"2020-03-18T11:00:00","end_utc":"2020-03-18T02:00:00+00:00","end_t":1584496800}]}'
j2=$(echo $j | jq ".result")
j3=$(echo $j2 | jq ".[0]")
printf "> "
if [ $r = "gachi" ]
then
printf ガチマッチ：
else
if [ $r = "league" ]
then
printf リーグマッチ：
else
printf レギュラーマッチ：
fi
fi
echo $j3 | jq -r ".rule"
j4=$(echo $j3 | jq ".maps_ex")
j5=$(echo $j4 | jq ".[0]")
j6=$(echo $j4 | jq ".[1]")
echo $j5 | jq -r ".image"
echo $j6 | jq -r ".image"
if [ $t = "now" ]
then
echo ">>> 現在のステージ："
else
echo ">>> 次のステージ："
fi
echo $j5 | jq -r ".name"
echo $j6 | jq -r ".name"
