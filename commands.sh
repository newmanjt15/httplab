#!/bin/bash

change_id=`date +%s`

res_dir="./tcp_logs/$change_id"

# mkdir -p $res_dir

log_file=$res_dir"/"$change_id".log"

cong_con="cubic"
acceptale_commands=(
    cubic
    bbr
    # TODO: reno
    )

if [ $# -eq 1 ]; then
    cong_con=$1
fi

if [[ " ${acceptale_commands[@]} " =~ " ${cong_con} " ]]; then
    if [ "$cong_con" == "cubic" ]; then
        echo `sudo cp ./cubic_sysctl.conf /etc/sysctl.conf`
    elif [ "$cong_con" == "bbr" ]; then
        echo `sudo cp ./bbr_sysctl.conf /etc/sysctl.conf`
    fi
    # (sudo sysctl --system > $log_file 2>&1)
    echo `sudo sysctl --system `
    new_cong=`sysctl net.ipv4.tcp_congestion_control | awk '{print $3;}'`
    echo "Your new congestion control is: $new_cong"
else
    echo "Please enter one of the following TCP congestion control commands: "
    echo "      - cubic"
    echo "      - bbr"
fi
