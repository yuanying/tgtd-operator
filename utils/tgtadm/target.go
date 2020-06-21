package tgtadm

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/longhorn/go-iscsi-helper/iscsi"
	longhornutil "github.com/longhorn/go-iscsi-helper/util"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
)

var (
	tgtBinary    = "tgtadm"
	targetPrefix = regexp.MustCompile(`^Target\s(\d+):\s+(.+)$`)
)

type TgtAdm interface {
	CreateTarget(tid int, name string) error
	DeleteTarget(tid int) error
	AddLunBackedByFile(tid int, lun int, backingFile string) error
	AddLun(tid int, lun int, backingFile string, bstype string, bsopts string) error
	DeleteLun(tid int, lun int) error
	BindInitiator(tid int, initiator string) error
	UnbindInitiator(tid int, initiator string) error
	GetTargetTid(name string) (int, error)
	GetTargetConnections(tid int) (map[string][]string, error)
	CloseConnection(tid int, sid, cid string) error
	FindNextAvailableTargetID() (int, error)

	GetTargets() ([]tgtdv1alpha1.TargetActual, error)
}

type TgtAdmLonghornHelper struct{}

func (t *TgtAdmLonghornHelper) CreateTarget(tid int, name string) error {
	return iscsi.CreateTarget(tid, name)
}
func (t *TgtAdmLonghornHelper) DeleteTarget(tid int) error {
	return iscsi.DeleteTarget(tid)
}
func (t *TgtAdmLonghornHelper) AddLunBackedByFile(tid int, lun int, backingFile string) error {
	return iscsi.AddLunBackedByFile(tid, lun, backingFile)
}
func (t *TgtAdmLonghornHelper) AddLun(tid int, lun int, backingFile string, bstype string, bsopts string) error {
	return iscsi.AddLun(tid, lun, backingFile, bstype, bsopts)
}
func (t *TgtAdmLonghornHelper) DeleteLun(tid int, lun int) error {
	return iscsi.DeleteLun(tid, lun)
}
func (t *TgtAdmLonghornHelper) BindInitiator(tid int, initiator string) error {
	return iscsi.BindInitiator(tid, initiator)
}
func (t *TgtAdmLonghornHelper) UnbindInitiator(tid int, initiator string) error {
	return iscsi.UnbindInitiator(tid, initiator)
}
func (t *TgtAdmLonghornHelper) GetTargetTid(name string) (int, error) {
	return iscsi.GetTargetTid(name)
}
func (t *TgtAdmLonghornHelper) GetTargetConnections(tid int) (map[string][]string, error) {
	return iscsi.GetTargetConnections(tid)
}
func (t *TgtAdmLonghornHelper) CloseConnection(tid int, sid, cid string) error {
	return iscsi.CloseConnection(tid, sid, cid)
}
func (t *TgtAdmLonghornHelper) FindNextAvailableTargetID() (int, error) {
	return iscsi.FindNextAvailableTargetID()
}
func (t *TgtAdmLonghornHelper) GetAccounts() ([]string, error) {
	opts := []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	output, err := longhornutil.Execute(tgtBinary, opts)
	if err != nil {
		return nil, err
	}
	return parseShowAccount(output)
}
func (t *TgtAdmLonghornHelper) GetTargets() (targets []tgtdv1alpha1.TargetActual, err error) {
	opts := []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	output, err := longhornutil.Execute(tgtBinary, opts)
	if err != nil {
		return nil, err
	}
	return parseShowTarget(output)
}

func parseShowAccount(raw string) ([]string, error) {
	success := false
	accounts := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Account list:") {
			success = true
		} else if len(line) > 0 {
			accounts = append(accounts, line)
		}
	}
	if !success {
		return nil, fmt.Errorf("Can't parse account information... not `Account list:` header")
	}
	return accounts, nil
}

func parseShowTarget(raw string) (targets []tgtdv1alpha1.TargetActual, err error) {
	targets = make([]tgtdv1alpha1.TargetActual, 0)
	scanner := bufio.NewScanner(strings.NewReader(raw))
	section := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		targetMatch := targetPrefix.FindStringSubmatch(line)
		if len(targetMatch) > 0 {
			section = ""
			target := tgtdv1alpha1.TargetActual{
				IQN: targetMatch[2],
			}
			tid, err := strconv.Atoi(strings.TrimSpace(targetMatch[1]))
			target.TID = int32(tid)
			if err != nil {
				return nil, fmt.Errorf("Cant parse TID from %s, line: %s", targetMatch[1], line)
			}
			targets = append(targets, target)
		} else if strings.HasPrefix(line, "LUN information:") {
			section = "lun"
		} else if strings.HasPrefix(line, "Account information:") {
			section = "account"
		} else if strings.HasPrefix(line, "ACL information:") {
			section = "acl"
		} else if section == "lun" {
			if strings.HasPrefix(line, "LUN:") {
				currentLUN := tgtdv1alpha1.TargetLUN{}
				lidstrs := strings.Split(line, ":")
				if len(lidstrs) != 2 {
					return nil, fmt.Errorf("Failed to parse and get LUN id: %v", line)
				}
				lid, err := strconv.Atoi(strings.TrimSpace(lidstrs[1]))
				if err != nil {
					return nil, err
				}
				currentLUN.LUN = int32(lid)
				target := currentTarget(targets)
				target.LUNs = append(target.LUNs, currentLUN)
			} else if strings.HasPrefix(line, "Backing store path:") {
				bspaths := strings.Split(line, ":")
				if len(bspaths) != 2 {
					return nil, fmt.Errorf("Failed to parse and get Backing store path: %v", line)
				}
				target := currentTarget(targets)
				bspath := strings.TrimSpace(bspaths[1])
				target.LUNs[len(target.LUNs)-1].BackingStore = bspath
			}
		} else if section == "account" {
			target := currentTarget(targets)
			target.Accounts = append(target.Accounts, line)
		} else if section == "acl" {
			target := currentTarget(targets)
			target.ACLs = append(target.ACLs, line)
		}
	}
	return
}

func currentTarget(targets []tgtdv1alpha1.TargetActual) *tgtdv1alpha1.TargetActual {
	return &targets[len(targets)-1]
}
