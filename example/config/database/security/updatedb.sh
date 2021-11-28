curl -L  "https://reputation.alienvault.com/reputation.snort" | sed '1,8d' | sed "s/ # /,/g" > tmp/reputation_alienvault.csv
curl -L "https://myip.ms/files/blacklist/general/latest_blacklist.txt" | sed '1,13d' | awk -F " " '{print $1",Malicious Host"}' > tmp/reputation_myip.csv
curl -L "https://data.netlab.360.com/feeds/hajime-scanner/bot.list" | awk -F "ip=" '{print $2",Botnet Scanner"}' > tmp/reputation_bot360.csv
cat tmp/reputation_alienvault.csv > reputation.csv
cat tmp/reputation_myip.csv >> reputation.csv
cat tmp/reputation_bot360.csv >> reputation.csv

#Botnet
curl -L "https://www.team-cymru.org/Services/Bogons/fullbogons-ipv4.txt" | sed '1d' | sed "/0.0.0.0/d" | sed "/224.0.0.0/d"  | sed "/240.0.0.0/d" | sed "/172.16.0.0/d" | sed "/192.168.0.0/d" | sed "/169.254.0.0/d" | sed "s/$/,BotNet/g" > tmp/botnet_cymru.txt
curl -L "https://www.spamhaus.org/drop/drop.txt" | sed '1,4d' | awk -F ' ' '{print $1}' | sed 's/$/,Spam Net/'  > tmp/botnet_spamhaus_drop.txt
curl -L "https://www.spamhaus.org/drop/edrop.txt" | sed '1,4d' | awk -F ' ' '{print $1}' | sed 's/$/,Spam Net/' > tmp/botnet_spamhaus_edrop.txt

cat tmp/botnet_cymru.txt > spam_bot_net.csv
cat tmp/botnet_spamhaus_drop.txt >> spam_bot_net.csv
cat tmp/botnet_spamhaus_edrop.txt >> spam_bot_net.csv
