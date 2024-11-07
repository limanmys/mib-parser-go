package mibparser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OIDNode represents a node in the MIB tree
type OIDNode struct {
	Name        string     `json:"name"`
	OID         string     `json:"oid"`
	ID          string     `json:"id"`
	Parent      string     `json:"parent"`
	Description string     `json:"description"`
	Children    []*OIDNode `json:"children,omitempty"`
}

// Parse reads and parses MIB files into OIDNodes
func (p *MIBParser) Parse() ([]OIDNode, error) {
	lines, err := p.ReadMIBFile()
	if err != nil {
		return nil, err
	}
	nodes, err := parseMIB(lines)
	if err != nil {
		return nil, err
	}
	return nodes, err
}

// GetJSONTree returns a JSON representation of the MIB tree
func (p *MIBParser) GetJSONTree() (string, error) {
	nodes, err := p.Parse()
	if err != nil {
		return "", err
	}
	tree, err := buildTree(nodes)
	if err != nil {
		return "", err
	}
	byteNodes, err := json.Marshal(tree)

	return string(byteNodes), nil
}

// GetObjects returns a JSON string of the parsed OIDNodes
func (p *MIBParser) GetObjects() (string, error) {
	nodes, err := p.Parse()
	if err != nil {
		return "", err
	}

	byteNodes, err := json.Marshal(nodes) // Marshal nodes to JSON
	if err != nil {
		return "", err
	}
	return string(byteNodes), nil
}

// buildTree constructs the MIB tree from OIDNodes
func buildTree(nodes []OIDNode) ([]*OIDNode, error) {
	nodeMap := make(map[string]*OIDNode)
	// Initialize the nodeMap
	for i := range nodes {
		node := &nodes[i]
		nodeMap[node.Name] = node
	}

	var rootNodes []*OIDNode
	rootNames := map[string]bool{"iso": true}
	for i := range nodes {
		node := &nodes[i]
		if rootNames[node.Parent] {
			rootNodes = append(rootNodes, node)
		} else if parent, exists := nodeMap[node.Parent]; exists {
			parent.Children = append(parent.Children, node)
		}
	}
	return rootNodes, nil
}

// parseMIB parses MIB lines into OIDNodes
func parseMIB(lines []string) ([]OIDNode, error) {
	var requiredMibs []string
	var definingMibs []string
	// Parse required and defining MIBs
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Contains(line, "FROM") {
			parts := strings.Split(line, "FROM")
			if len(parts) >= 2 {
				subParts := strings.Split(parts[1], ";")
				if len(subParts) > 0 {
					requiredMib := strings.TrimSpace(subParts[0])
					requiredMibs = append(requiredMibs, requiredMib+".mib")
				}
			}
		}
		if strings.Contains(line, "DEFINITIONS ::= BEGIN") {
			parts := strings.Split(line, "DEFINITIONS ::= BEGIN")
			if len(parts) > 0 {
				definingMib := strings.TrimSpace(parts[0])
				definingMibs = append(definingMibs, definingMib+".mib")
			}
		}
	}
	// Remove redundant required MIBs
	for i := len(requiredMibs) - 1; i >= 0; i-- {
		for _, definingMib := range definingMibs {
			if definingMib == requiredMibs[i] {
				requiredMibs = append(requiredMibs[:i], requiredMibs[i+1:]...)
				break
			}
		}
	}
	uniqueMap := make(map[string]bool)
	uniqueReqireds := []string{}
	for _, item := range requiredMibs {
		if _, found := uniqueMap[item]; !found {
			uniqueMap[item] = true
			uniqueReqireds = append(uniqueReqireds, item)
		}
	}
	if len(uniqueReqireds) > 0 {
		return nil, fmt.Errorf("parsing operation has failed :\n \t\t\t\t\trequired files: %s", uniqueReqireds)
	}
	var nodes []OIDNode
	// Parse lines into OIDNodes
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Contains(line, "OBJECT IDENTIFIER ::= {") {
			parts := strings.Split(line, "OBJECT IDENTIFIER ::= {")
			if len(parts) >= 2 {
				nameParts := strings.Split(parts[0], " ")
				nextNodeParts := strings.Split(strings.TrimSpace(strings.Trim(parts[1], "}")), " ")
				if len(nameParts) > 0 && len(nextNodeParts) >= 2 {
					name := strings.TrimSpace(nameParts[0])
					nodes = append(nodes, OIDNode{
						Name:        name,
						ID:          nextNodeParts[1],
						Parent:      nextNodeParts[0],
						Description: "",
					})
				}
			}
		} else if (strings.Contains(line, "OBJECT-TYPE") || strings.Contains(line, "OBJECT-IDENTITY")) && !strings.Contains(line, "MODULE-IDENTITY") {
			nameParts := strings.Fields(line)
			if len(nameParts) > 0 {
				name := strings.TrimSpace(nameParts[0])
				parent := ""
				description := ""
				for j := i + 1; j < len(lines); j++ {
					nextLine := strings.TrimSpace(lines[j])
					if strings.HasPrefix(nextLine, "::= {") {
						parentParts := strings.Split(nextLine, "{")
						if len(parentParts) > 1 {
							parent = strings.Trim(parentParts[1], " }")
						}
						break
					}
					if strings.HasPrefix(nextLine, "DESCRIPTION") {
						descriptionLine := nextLine
						for k := j + 1; k < len(lines); k++ {
							if strings.Contains(lines[k], "::= {") {
								break
							}
							descriptionLine += " " + strings.TrimSpace(lines[k])
						}
						description = strings.TrimSpace(strings.ReplaceAll(description, "DESCRIPTION", ""))
					}
				}
				if !strings.Contains(name, "OBJECT-TYPE") && !strings.Contains(name, "--") && !strings.Contains(name, "OBJECT-IDENTITY") && !strings.Contains(line, "MODULE-IDENTITY") {
					parentParts := strings.Split(parent, " ")
					if len(parentParts) >= 2 {
						nodes = append(nodes, OIDNode{
							Name:        name,
							ID:          parentParts[1],
							Parent:      parentParts[0],
							Description: description,
						})
					}
				}
			}
		}
	}
	formatedNodes, err := setOids(nodes) // Set OIDs for nodes
	if err != nil {
		return nil, err
	}
	return formatedNodes, nil
}

func setOids(nodes []OIDNode) ([]OIDNode, error) {

	// Initialize formatted nodes slice
	var formatedNodes []OIDNode

	// Create a map to store nodes by name
	nodeMap := make(map[string]OIDNode)
	for _, node := range nodes {
		nodeMap[node.Name] = node
	}

	// Iterate over nodes to set OIDs
	for _, node := range nodes {
		oidParts := []string{node.ID} // Start with current node ID
		parent := node.Parent

		for parent != "" {
			if nextNode, found := nodeMap[parent]; found {
				oidParts = append([]string{nextNode.ID}, oidParts...)
				parent = nextNode.Parent

			} else {
				break
			}
		}

		oid := strings.Join(oidParts, ".") // Join all parts to form full OID
		formatedNodes = append(formatedNodes, OIDNode{
			Name:        node.Name,
			ID:          node.ID,
			Parent:      node.Parent,
			OID:         "1." + oid,
			Description: node.Description,
		})
	}

	return formatedNodes, nil
}
