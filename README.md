tgtd-operator
=============

![envtest](https://github.com/yuanying/tgtd-operator/workflows/envtest/badge.svg)

tgtd-operator is a light-weight Kubernetes Operator that operates tgt daemon to manage iSCSI Target, LUN, Account, Initiator.

Architecture
------------

tgtd-operator will communicate with tgt daemon using tgtadm cli, so tgtd-operator should be placed on each nodes where tgtd is running. And tgtd-operator uses Target/InitiatorGroup/InitiatorGroupBinding/Account CRD which are stored in Kubernetes apiserver instead of using `/etc/tgt/targets.conf`.

Concepts
--------

### Target

Target is a representation of iSCSI Target which consits of LUNs.

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

InitiatorGroup is a group of iSCSI initiators. It has nodeSelector to create initiator group from Node objects and so it means InitiatorGroup dynamically discover initiator groups.
Usually, it should not to set `nodeSelector` because all nodes want to be initiator.

```yaml
apiVersion: tgtd.unstable.cloud/v1alpha1
kind: InitiatorGroup
metadata:
  name: initiator-group-sample
spec:
  addresses: ["192.168.1.0/24", "192.168.2.0/24"]
  nodeSelector:
    "initiator-group": "group1"
  initiatorNameStrategy:
    type: NodeName
    initiatorNamePrefix: "iqn.2020-06.cloud.unstable.tgtd.init"
```

"initiator name" is decided by InitiatorGroup from node name or annotation.
If you specify `.spec.initiatorNameStrategy.type` to `NodeName` and `spec.initiatorNameStrategy.initiatorNamePrefix` to `iqn.2020-06.cloud.unstable.tgtd.init`, "initiator name" is generated from node name and then it become `iqn.2020-06.cloud.unstable.tgtd.init:${NODE_NAME}`.
Or, if you specify `.spec.initiatorNameStrategy.type` to `Annotation` and `.spec.initiatorNameStrategy.AnnotationKey`, "initiator name" is decided to it's annotation value.

### InitiatorGroupBinding

InitiatorGroupBinding binds a Target with a InitiatorGroup directly. If the binding succeeds, tgtd-operator tries to bind target with initiator using `tgtadm` utility.

Sample configuration is below.

```yaml
apiVersion: tgtd.unstable.cloud/v1alpha1
kind: InitiatorGroupBinding
metadata:
  name: initiator-group-binding-sample
spec:
  targetRef: target-sample
  initiatorGroupRef: initiator-group-sample
```
