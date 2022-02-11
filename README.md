# golang-yaml-rules/yaml-transform
Library/tool to change a yaml given a rules file

Using jsonpath ( https://github.com/vmware-labs/yaml-jsonpath ), this tool shows how to 
transform a yaml source using simple rules

- 'set' rule:  set/update values (single value, or whole yaml section) within nodes that match a first-level matching JsonPath expression
       values need not be applied to the matching node itself : a 'subpath' jsonpath expression allows to chose on which subnode you want to 
       set values.
       
- 'deleteChildrenThatMatch' rule:  removes array children or keys from nodes that match a first-level matching JsonPath expression.
       
       When applied to an Array, the children to remove are those matching a second-level matching JsonPath (syntax following arrays conditional within the '[]' JsonPath construct)
       
       When applied to a Map Object, this also allows to remove a single key by its name. But in this case, wildcard or conditions based on the object content are not supported (only a single fixed key is removed)


As an example, see "rules_example.yaml", and "myfile_to_transform.yaml"


Usage:

       yaml-transform <rules yaml file path> <document yaml file path>

try out:
       
       ./yaml-transform  rules_example.yaml myfile_to_transform.yaml

This will output transformed yaml file to stdout

