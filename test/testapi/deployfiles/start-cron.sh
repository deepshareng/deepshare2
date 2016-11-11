#!/bin/sh
# start-cron.sh

# resolve the cron unrunnable issue by following this solution :
# http://raid6.com.au/posts/docker-cron-rsyslog/
rsyslogd
touch /var/log/cron.log
touch /etc/crontab /etc/cron.*/*
cron
tail -F /var/log/syslog /var/log/cron.log
