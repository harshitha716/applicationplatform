package models

import "github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"

type NodeType string

const (
	NodeTypeDataset NodeType = "dataset"
	NodeTypeFolder  NodeType = "folder"
)

type DAGNode struct {
	NodeId     string
	NodeType   NodeType
	Parents    []*DAGNode
	EdgeConfig map[string]interface{}
}

func (d *DAGNode) GetImportFilePath() (string, error) {
	if d.NodeType == NodeTypeFolder {
		return d.NodeId, nil
	}

	if len(d.Parents) == 1 {
		return d.Parents[0].GetImportFilePath()
	}

	return "", errors.ErrDatasetNotEligibleForFileImport
}
