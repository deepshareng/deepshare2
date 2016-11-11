#!/bin/bash

test_result=`/go/bin/deepshare2_apitest.test -env staging -logfile /opt/deepshare2-api-error2.log -genurl-addr=http://deepshare2-genurl.ds-production:16759 -inappdata-addr=http://deepshare2-inappdata.ds-production:16759 -binddevicetocookie-addr=http://deepshare2-bind.ds-production:16759 -sharelink-addr=http://deepshare2-sharelink.ds-production:16759 -dsusage-addr=http://deepshare2-dsusage.ds-production:16759`
pass=`echo $test_result | grep PASS`
error_info=`cat /opt/deepshare2-api-error2.log`


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
