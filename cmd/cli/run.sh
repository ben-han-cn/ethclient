gethpath=`which geth`
killall geth
rm -rf keys/* && rm -rf nodedata/* && ./cli -k keys -n nodedata -p $gethpath -s 1 -y 1
