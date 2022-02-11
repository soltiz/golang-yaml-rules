package main

import (
	"bytes"
	"fmt"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func readInput() yaml.Node {
	yamlFile, err := ioutil.ReadFile("mycrd.yaml")

	if err != nil {
		log.Fatal("Unable to read input file", err)
	}

	var rootNode yaml.Node

	err = yaml.Unmarshal([]byte(yamlFile), &rootNode)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	return rootNode
}

type SetSpec struct {
	Subpath string
	Values  map[string]yaml.Node
}

type RuleSpec struct {
	Match                   string
	Set                     []SetSpec
	DeleteChildrenThatMatch string `yaml:"deleteChildrenThatMatch"`
}

func readRules() map[string]RuleSpec {
	yamlFile, err := ioutil.ReadFile("rules.yaml")

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]RuleSpec)
	err2 := yaml.Unmarshal(yamlFile, &data)

	if err2 != nil {
		log.Fatal("Unable to unmarshall rules yaml file", err2)
	}
	return data
}

func applyRule(ruleSpec RuleSpec, document yaml.Node) {
	p, err := yamlpath.NewPath(ruleSpec.Match)
	if err != nil {
		log.Fatalf("cannot create path: %v", err)
	}

	q, err := p.Find(&document)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	for _, matchingNode := range q {
		if ruleSpec.Set != nil {
			for _, setSpec := range ruleSpec.Set {
				applySet(setSpec, matchingNode)
			}
		} else {
			deleteChildrenThatMatch(ruleSpec.DeleteChildrenThatMatch, matchingNode)
		}
	}
}

func deleteChildrenThatMatch(match string, parentNode *yaml.Node) {
	childrenMatchString := "$[" + match + "]"
	p, err := yamlpath.NewPath(childrenMatchString)
	if err != nil {
		log.Fatalf("cannot create children match lookup: %v", err)
	}

	nodesToRemove, err := p.Find(parentNode)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	nodesToKeep := []*yaml.Node{}

	for _, childNode := range parentNode.Content {
		if !nodesInList(childNode, nodesToRemove) {
			nodesToKeep = append(nodesToKeep, childNode)
		}
	}

	parentNode.Content = nodesToKeep
}

func nodesInList(node *yaml.Node, nodesList []*yaml.Node) bool {
	for _, item := range nodesList {
		if item == node {
			return true
		}
	}
	return false
}

func buildStringNodes(key string, value string) []*yaml.Node {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
	return []*yaml.Node{keyNode, valueNode}
}

func buildKeyNode(key string) *yaml.Node {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}

	return keyNode
}

func applySet(spec SetSpec, node *yaml.Node) {

	p, err := yamlpath.NewPath(spec.Subpath)
	if err != nil {
		log.Fatalf("cannot create subpath: %v", err)
	}

	matches, err := p.Find(node)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	for _, matchingNode := range matches {
		for fieldName, fieldValue := range spec.Values {

			fieldPath, err := yamlpath.NewPath("$." + fieldName)
			if err != nil {
				log.Fatalf("cannot create field subpath: %v", err)
			}

			fieldMatches, err := fieldPath.Find(matchingNode)

			if len(fieldMatches) != 0 {
				r := &fieldMatches[0]
				(*r).Kind = fieldValue.Kind
				(*r).Tag = fieldValue.Tag
				(*r).Value = fieldValue.Value
				(*r).Content = fieldValue.Content
				(*r).Column = fieldValue.Column
			} else {
				matchingNode.Content = append(matchingNode.Content, buildKeyNode(fieldName), &fieldValue)
			}
		}
	}
}

func outputResult(doc yaml.Node) {
	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)
	defer e.Close()
	e.SetIndent(2)

	err := e.Encode(&doc)
	if err != nil {
		log.Printf("Error: cannot marshal node: %v", err)
		return
	}
	fmt.Println(buf.String())
}

func main() {
	rules := readRules()
	document := readInput()
	for ruleName, ruleSpec := range rules {
		log.Printf("Applying rule '%v'...\n", ruleName)
		applyRule(ruleSpec, document)
	}
	outputResult(document)
}
