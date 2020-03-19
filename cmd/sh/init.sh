#!/bin/bash
#init script for using container

. mashironrc
rm -r $VM
mkdir $VM
pacstrap -i -c $VM bash busybox
systemd-nspawn -D $VM /usr/bin/busybox --install

#container.sh
#$script = "#!/bin/bash
#... script ..."
#$option="command line options"
cat > $VM/usr/container.sh << 'EOF'
#!/bin/bash
#container script
. <(echo "$2")
timeout -sKILL 3 bash <(echo "$1") ${@:3}
STATUS=$?
if [ "$STATUS" -eq 124 ];
then
        echo "Error:timeout"
        exit 124
fi
EOF
chmod +x $VM/usr/container.sh
