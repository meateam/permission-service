#!/bin/bash

for i in {1..1000}
do
    echo "hello $i" 
    content="{  \"fileID\": \"f${i}\",  \"userID\": \"u1\",  \"role\": 1,  \"creator\": \"u2\",  \"override\": true,  \"appID\": \"drive\"}"
    echo $content
    /home/shahar/go/bin/grpcurl -plaintext -d "$content" 0.0.0.0:8087 permission.Permission/CreatePermission
done

# grpcurl -plaintext -d '{  "userID": "u1",  "pageNum": 0,  "pageSize": 20,  "isShared": true,  "appID": "drive"}' 0.0.0.0:8087 permission.Permission/GetUserPermissions
