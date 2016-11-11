#!/bin/bash

test_result=`/go/bin/deepshare2_apitest.test -env production -logfile /opt/deepshare2-api-error1.log`
pass=`echo $test_result | grep PASS`
error_info=`cat /opt/deepshare2-api-error1.log`


if [ "$pass" != "PASS" ]; then
    test_result=`echo $test_result`
    error_info=${error_info//'\n'//'<br>'}
    echo -e "SUBJECT:$test_result
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="XX-1234DED00099A";
Content-Transfer-Encoding: 7bit

--XX-1234DED00099A
Content-Type: text/plain;charset=UTF-8
Content-Transfer-Encoding: 7bit

$error_info" | msmtp -t deepshare-alert@misingularity.com
fi
