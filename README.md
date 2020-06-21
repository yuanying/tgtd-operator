tgtd-operator
=============

tgtd-operator is a light-weight Kubernetes Operator that operates tgt daemon to manage iSCSI Target, LUN, Account, Initiator.

Architecture
------------

tgtd-operator will communicate with tgt daemon using tgtadm cli, so tgtd-operator should be placed on each nodes where tgtd is running. And tgtd-operator uses Target/InitiatorGroup/InitiatorGroupBinding/Account CRD which are stored in Kubernetes apiserver instead of using `/etc/tgt/targets.conf`.

### Target

Target is a represent of iSCSI Target which consits of LUNs.

```yaml
apiVersion: tgtd.unstable.cloud/v1alpha1
kind: Target
metadata:
  name: target-sample
spec:
  // nodeName is a node name where the target will be placed.
  nodeName: storage-node
  // iqn is an iqn of the target
  iqn: iqn.2020-04.cloud.unstable.example:target
  // luns is a list of LUNs
  luns:
    // lun is an id of the LUN
  - lun: 1
    // backingStore is a path of the backing store
    backingStore: /dev/sda1
  - lun: 2
    backingStore: /dev/sdb1
    bsType: aio
```

### InitiatorGroup

### InitiatorGroupBinding
