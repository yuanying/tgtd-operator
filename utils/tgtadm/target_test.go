package tgtadm

import (
	"reflect"
	"testing"
)

var testAccountsStr = `
Account list:
    user3
    user2
    user1
`

var testTargetsStr = `
Target 1: iqn.2020-04.cloud.unstable:target1
    System information:
        Driver: iscsi
        State: ready
    I_T nexus information:
    LUN information:
        LUN: 0
            Type: controller
            SCSI ID: IET     00010000
            SCSI SN: beaf10
            Size: 0 MB, Block size: 1
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: null
            Backing store path: None
            Backing store flags:
        LUN: 1
            Type: disk
            SCSI ID: IET     00010001
            SCSI SN: beaf11
            Size: 5369 MB, Block size: 512
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: rdwr
            Backing store path: /dev/zvol/tank/vol
            Backing store flags:
        LUN: 3
            Type: disk
            SCSI ID: IET     00010003
            SCSI SN: beaf13
            Size: 5369 MB, Block size: 512
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: rdwr
            Backing store path: /dev/zvol/tank/vol2
            Backing store flags:
        LUN: 13
            Type: disk
            SCSI ID: IET     0001000d
            SCSI SN: beaf113
            Size: 5369 MB, Block size: 512
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: rdwr
            Backing store path: /dev/zvol/tank/vol3
            Backing store flags:
    Account information:
        user1
    ACL information:
        ALL
Target 3: iqn.2020-04.cloud.unstable:target3
    System information:
        Driver: iscsi
        State: ready
    I_T nexus information:
    LUN information:
        LUN: 0
            Type: controller
            SCSI ID: IET     00030000
            SCSI SN: beaf30
            Size: 0 MB, Block size: 1
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: null
            Backing store path: None
            Backing store flags:
        LUN: 1
            Type: disk
            SCSI ID: IET     00030001
            SCSI SN: beaf31
            Size: 5369 MB, Block size: 512
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: rdwr
            Backing store path: /dev/zvol/tank/vol
            Backing store flags:
    Account information:
        user1
        user2
        user3
    ACL information:
        192.168.1.0/24
        192.168.2.0/24
        192.168.3.0/24
        iqn.2020-04.cloud.unstable:init1
Target 4: iqn.2020-04.cloud.unstable:target4
    System information:
        Driver: iscsi
        State: ready
    I_T nexus information:
    LUN information:
        LUN: 0
            Type: controller
            SCSI ID: IET     00040000
            SCSI SN: beaf40
            Size: 0 MB, Block size: 1
            Online: Yes
            Removable media: No
            Prevent removal: No
            Readonly: No
            SWP: No
            Thin-provisioning: No
            Backing store type: null
            Backing store path: None
            Backing store flags:
    Account information:
    ACL information:
`

func TestParseTargetsOutput(t *testing.T) {

	got, err := parseShowTarget(testTargetsStr)
	if err != nil {
		t.Fatalf("parseShowTarget should be success: %v", err)
	}
	want := []Target{
		{
			TID: 1,
			IQN: "iqn.2020-04.cloud.unstable:target1",
			LUNs: []LUN{
				{
					ID:               0,
					BackingStorePath: "None",
				},
				{
					ID:               1,
					BackingStorePath: "/dev/zvol/tank/vol",
				},
				{
					ID:               3,
					BackingStorePath: "/dev/zvol/tank/vol2",
				},
				{
					ID:               13,
					BackingStorePath: "/dev/zvol/tank/vol3",
				},
			},
			Accounts: []string{"user1"},
			ACLs:     []string{"ALL"},
		},
		{
			TID: 3,
			IQN: "iqn.2020-04.cloud.unstable:target3",
			LUNs: []LUN{
				{
					ID:               0,
					BackingStorePath: "None",
				},
				{
					ID:               1,
					BackingStorePath: "/dev/zvol/tank/vol",
				},
			},
			Accounts: []string{"user1", "user2", "user3"},
			ACLs:     []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "iqn.2020-04.cloud.unstable:init1"},
		},
		{
			TID: 4,
			IQN: "iqn.2020-04.cloud.unstable:target4",
			LUNs: []LUN{
				{
					ID:               0,
					BackingStorePath: "None",
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Result is not equal expected, \ngot: %v\nwant: %v", got, want)
	}
}

func TestParseShowAccount(t *testing.T) {
	got, err := parseShowAccount(testAccountsStr)
	if err != nil {
		t.Fatalf("parseShowAccount should be success: %v", err)
	}
	want := []string{"user3", "user2", "user1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Result is not equal expected, \ngot: %v\nwant: %v", got, want)
	}
}
