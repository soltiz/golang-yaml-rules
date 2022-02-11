package main

import (
	"bytes"
	"fmt"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"log"
)

func Example() {
	y := `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
      - name: nginy
        image: nginy
        ports:
        - containerPort: 81
`
	var n yaml.Node

	err := yaml.Unmarshal([]byte(y), &n)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	p, err := yamlpath.NewPath("$..spec.containers[*].image")
	if err != nil {
		log.Fatalf("cannot create path: %v", err)
	}

	q, err := p.Find(&n)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	for _, i := range q {
		i.Value = "example.com/user/" + i.Value
	}

	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)
	defer e.Close()
	e.SetIndent(2)

	err = e.Encode(&n)
	if err != nil {
		log.Printf("Error: cannot marshal node: %v", err)
		return
	}

	fmt.Println(buf.String())

}

func main() {
	Example()
	log.Fatalf("Unable to read client secret file")
}
