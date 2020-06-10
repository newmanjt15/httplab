#!/bin/bash

# parameters
test_id=`date +%s`
res_folder="./wget_results/$test_id"

if [ $# -eq 1 ]
then
    if [[ "$1" == *"https://"* ]] || [[ "$1" == *"http://"* ]]
    then
        # echo "$1" > $url_file
        base_url=$1
    else
        url_file=$1
    fi
elif [ $# -eq 2 ]
then
    if [[ "$1" == *"https://"* ]] || [[ "$1" == *"http://"* ]]
    then
        # echo "$1" > $url_file
        base_url=$1
    else
        url_file=$1
    fi
    res_folder=$2
fi
url_file="./.url"
depth=0

mkdir -p $res_folder
chmod -R 777 $res_folder



# echo "$base_url"

upload_page () {
    url=$1
    domain=`echo $url | awk -F "/" '{print $3}' | awk -F "." '{print $1;}'`
    directory=$res_folder"/"$domain
    scp -q -r -i ~/.ssh/reporter.pem $directory ubuntu@ec2-54-91-253-3.compute-1.amazonaws.com:/var/www/testmyprotocol.com/public/cloned_sites/
    link="https://testmyprotocol.com/cloned_sites/"$domain"/"
}

download_page () {
    url=$1
    domain=`echo $url | awk -F "/" '{print $3}' | awk -F "." '{print $1;}'`
    # echo "Downloading $url from $domain"
    mkdir -p $res_folder"/"$domain
    wget_flags=(
        -E
        -H
        -k
        -K
        -p
        -q
        --random-wait 
        --ignore-length 
        --header="User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
        # --header="Accept-Encoding: gzip;q=0"
        --no-check-certificate
    )
    wget "${wget_flags[@]}" -P $res_folder"/"$domain $url
    name=`echo "$url" | awk -F "/" '{print $3}'`
    chmod -R 755 $res_folder"/"$domain 
    #handle compression; some servers will send compressed even if we don't want it
    compressed_files=`find $res_folder"/"$domain"/"`    
    index_file="$res_folder/$domain/$name/index.html"
    while IFS= read compressed_file
    do
        if [[ "$compressed_file" == *".gz"* ]]
        then
            tempname1=`echo "$compressed_file" | awk -F "/" '{print $NF}'`
            if [ -f "$index_file" ]; then
                sed -i -e "s~$tempname1~${tempname1%.gz}~" $index_file
            fi
            gzip -d $compressed_file
        fi
    done <<< "$compressed_files"
    name=`echo "$url" | awk -F "/" '{print $3}'`
    # insert our script into the index page to track metrics
    if [ -f "$index_file" ]; then
        sed -i -e "s~<\/body>~<script src='\/cloned_sites\/js\/insert_script.js'><\/script><\/body>~" $index_file
    fi
    echo "$res_folder/$domain/$name"
}

log_file=$res_folder"/"$test_id".test_log"

(download_page "$base_url" > $log_file 2>&1)
# download_page "$base_url"

echo `cat $log_file`
