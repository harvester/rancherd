package controller

import (
	"testing"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_NodeElementalUpgrade(t *testing.T) {
	const nodeName = "test-node"
	var testCases = []struct {
		name        string
		nodeObj     *corev1.Node
		expectedErr bool
	}{
		{
			name: "test upgrade success",
			nodeObj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
					Annotations: map[string]string{
						ElementalUpgradeRequestAnnotationKey:    ElementalUpgradeRequestAnnotationValue,
						ElementalOverrideCommandAnnotationKey:   "echo",
						ElementalOverrideArgumentsAnnotationKey: "successful upgrade",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "test upgrade failure",
			nodeObj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
					Annotations: map[string]string{
						ElementalUpgradeRequestAnnotationKey:    ElementalUpgradeRequestAnnotationValue,
						ElementalOverrideCommandAnnotationKey:   "false",
						ElementalOverrideArgumentsAnnotationKey: "",
					},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := require.New(t)
			client := fake.NewSimpleClientset(tc.nodeObj)
			m := &ElementalUpgradeManager{
				nodeName:       nodeName,
				nodeController: client.CoreV1().Nodes(),
			}

			nodeObj, err := m.OnChange(nodeName, tc.nodeObj)
			// expect no error during actual error, the failure of upgrade
			// is reflect in condition of the node object
			assert.NoError(err, "unexpected error during reconcile of object")
			elementalUpgradeCondition := getElementalCondition(nodeObj)
			assert.NotNil(elementalUpgradeCondition, "elemental upgrade condition should not be nil")
			t.Log(elementalUpgradeCondition)
			if tc.expectedErr {
				assert.Equal(corev1.ConditionFalse, elementalUpgradeCondition.Status, "expected upgrade to fail")
			} else {
				assert.Equal(corev1.ConditionTrue, elementalUpgradeCondition.Status, "expected upgrade to succeed")
			}
		})
	}
}

func getElementalCondition(nodeObj *corev1.Node) *corev1.NodeCondition {
	for _, condition := range nodeObj.Status.Conditions {
		if condition.Type == ElementalUpgradeStatusCondition {
			return &condition
		}
	}
	return nil
}
