#!/bin/bash

#change maximum File-Descriptors
echo "fs.file-max = 100000" >> /etc/sysctl.conf

#change maximum concurrency
echo "net.core.somaxconn = 65536" >> /etc/sysctl.conf


echo "* soft nproc 65535" >> /etc/security/limits.conf
echo "* hard nproc 65535" >> /etc/security/limits.conf
echo "* soft nofile 65535" >> /etc/security/limits.conf
echo "* hard nofile 65535" >> /etc/security/limits.conf
