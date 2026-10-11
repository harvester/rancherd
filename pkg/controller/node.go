package controller

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ctrlv1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	nodev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	ElementalUpgradeRequestAnnotationKey    = "harvesterhci.io/elemental-upgrade-requested"
	ElementalUpgradeRequestAnnotationValue  = "true"
	ElementalUpgradeStatusCondition         = "ElementalUpgradeStatus"
	ElementalUpgradeArtifacts               = "/tmp/elemental-upgrade-artifacts"
	ElementalUpgradeArtifactsAnnotationKey  = "harvesterhci.io/elemental-upgrade-artifacts"
	elementalBinaryName                     = "elemental"
	defaultElementalPath                    = "/usr/bin"
	configDirectory                         = "config"
	rootfsDirectory                         = "rootfs"
	configFileName                          = "config.yaml"
	rancherdConfigDirectory                 = "/etc/rancher/rancherd"
	ElementalOverrideCommandAnnotationKey   = "harvesterhci.io/elemental-override-command"
	ElementalOverrideArgumentsAnnotationKey = "harvesterhci.io/elemental-override-arguments"
	ElementalUpgradeSuccessfulReason        = "ElementalUpgradeSuccessful"
	ElementalUpgradeFailedReason            = "ElementalUpgradeFailed"
)

type ElementalUpgradeManager struct {
	ctx            context.Context
	nodeController nodev1.NodeInterface
	nodeName       string
}

func Register(ctx context.Context, nodeController ctrlv1.NodeController, nodeInterface nodev1.NodeInterface) error {
	var err error
	nodeName, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	mgr := &ElementalUpgradeManager{
		ctx:            ctx,
		nodeController: nodeInterface,
		nodeName:       nodeName,
	}
	nodeController.OnChange(ctx, "rancherd-node-controller", mgr.OnChange)
	return nil
}

func (m *ElementalUpgradeManager) OnChange(key string, nodeObj *corev1.Node) (*corev1.Node, error) {
	if key != m.nodeName {
		logrus.Debugf("node key %s does not match hostname %s", key, m.nodeName)
		return nodeObj, nil
	}

	if nodeObj == nil || nodeObj.DeletionTimestamp != nil || nodeObj.Annotations[ElementalUpgradeRequestAnnotationKey] != ElementalUpgradeRequestAnnotationValue {
		return nodeObj, nil
	}

	// perform elemental upgrade if needed
	return m.runElementalUpgrade(nodeObj)
}

// runElementalUpgrade executes the elemental upgrade process for the given node using the specified artifacts and configuration.
func (m *ElementalUpgradeManager) runElementalUpgrade(nodeObj *corev1.Node) (*corev1.Node, error) {
	// to ensure upgrade is only run on a node where it was successful, we need to ensure that elemental upgrade condition
	// has not been set on the node
	exists := checkElementalUpgradeConditionExists(nodeObj)
	if exists {
		logrus.Infof("elemental upgrade has already been attempted on the node, skipping node %s", m.nodeName)
		return nodeObj, nil
	}
	artifactPath := ElementalUpgradeArtifacts
	elementalBinary := filepath.Join(defaultElementalPath, elementalBinaryName)
	val, ok := nodeObj.Annotations[ElementalUpgradeArtifactsAnnotationKey]
	if ok && val != "" {
		artifactPath = val
	}
	_, err := os.Stat(filepath.Join(artifactPath, elementalBinaryName))
	if err == nil {
		elementalBinary = filepath.Join(artifactPath, elementalBinaryName)

	}

	// default config directory is artifactPath/config
	configDir := filepath.Join(artifactPath, configDirectory)
	// check if override elemental upgrade config exists at /etc/rancher/rancherd/upgrade.yaml
	upgradeConfigPath := filepath.Join(rancherdConfigDirectory, "upgrade.yaml")
	if _, err := os.Stat(upgradeConfigPath); err == nil {
		configDir = rancherdConfigDirectory
	}

	rootfsDir := filepath.Join(artifactPath, rootfsDirectory)

	elementalArgs := []string{
		"--logfile", filepath.Join(artifactPath, "elemental-upgrade.log"),
		"--directory", rootfsDir,
		"--config-dir", configDir,
		"--debug",
	}
	// check for elemental override command and arguments. This makes it easier to customize the upgrade process per node.
	// this can also be used to perform unit tests to verify reconcile of nodes and conditions
	if val, ok := nodeObj.Annotations[ElementalOverrideCommandAnnotationKey]; ok && val != "" {
		elementalBinary = val
	}
	if val, ok := nodeObj.Annotations[ElementalOverrideArgumentsAnnotationKey]; ok && val != "" {
		elementalArgs = strings.Fields(val)
	}

	logrus.Debugf("Running elemental upgrade with binary: %s and args: %v\n", elementalBinary, elementalArgs)
	cmd := exec.Command(elementalBinary, elementalArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("failed to run elemental upgrade on node %s: %v, output: %s", m.nodeName, err, string(output))
	}
	logrus.Debugf("elemental upgrade output: %s", string(output))
	return m.updateElementalUpgradeCondition(nodeObj, err)
}

// check if ElementalUpgradeStatusCondition exists and return its value
func checkElementalUpgradeConditionExists(nodeObj *corev1.Node) bool {
	for _, condition := range nodeObj.Status.Conditions {
		if condition.Type == ElementalUpgradeStatusCondition {
			return true
		}
	}
	return false
}

func (m *ElementalUpgradeManager) updateElementalUpgradeCondition(nodeObj *corev1.Node, err error) (*corev1.Node, error) {
	condition := corev1.NodeCondition{
		Type:               ElementalUpgradeStatusCondition,
		LastHeartbeatTime:  metav1.Now(),
		LastTransitionTime: metav1.Now(),
	}
	if err == nil {
		condition.Status = corev1.ConditionTrue
		condition.Reason = ElementalUpgradeSuccessfulReason
		condition.Message = "Elemental upgrade succeeded"
	} else {
		condition.Status = corev1.ConditionFalse
		condition.Reason = ElementalUpgradeFailedReason
		condition.Message = err.Error()
	}
	nodeObj.Status.Conditions = append(nodeObj.Status.Conditions, condition)
	return m.nodeController.UpdateStatus(m.ctx, nodeObj, metav1.UpdateOptions{})
}
