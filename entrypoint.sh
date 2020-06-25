#!/bin/bash

cleanup() {
    echo "Cleanup tgtd..."
    tgtadm --op update --mode sys --name State -v offline
    tgt-admin --offline ALL
    tgt-admin --update ALL -c /dev/null -f
    tgtadm  --op delete --mode system
}
trap 'cleanup' SIGTERM

tgtd -f &
tgtadm --op update --mode sys --name State -v offline
tgt-admin -e -c /etc/tgt/targets.conf
tgtadm --op update --mode sys --name State -v ready

wait $!
