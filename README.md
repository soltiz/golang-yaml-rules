# golang-yaml-rules/yaml-transform
Library/tool to change a yaml given a rules file

Using jsonpath ( https://github.com/vmware-labs/yaml-jsonpath ), this tool shows how to 
transform a yaml source using simple rules

The rules file is a yaml map, the keys being only names.


- 'set' rule:  set/update values (single value, or whole yaml section) within nodes that match a first-level matching JsonPath expression
       values need not be applied to the matching node itself : a 'subpath' jsonpath expression allows to chose on which subnode you want to 
       set values.

  ```yaml
       # Syntax:
       
       - name_of_rule: whatever
         # match is a node selector, starting at your document root ($)

         # For some hints on Jsonpath syntax, you can have a look for example at https://support.smartbear.com/alertsite/docs/monitors/api/endpoint/jsonpath.html

         #  Here the selector targets all items in the 'containers' array
         match: "$..spec.containers[*]" 
                set:
                 # Here you may have multiple groups of values you want to apply
                 # for each, you select the subnodes of your matching node (so here, '$' is not the root
                 # of your document, but the yaml node that matched the 'match' jsonpath expression)
                - subpath: $
                  # Here you provide IMPERATIVE values that will override the potential ones with same name
                  # in your selected subnodes . This values can be simple types or multi-level yaml objects
                  values:
                    someField: a value
                    someOtherField:
                       a_subfield: true
                       other_subfield: 42
                - subpath: $.some_items_collection[*]
                  values:
                     all_my_subitems_get_an:  A

  ```



- 'deleteChildrenThatMatch' rule:  removes array children or keys from nodes that match a first-level matching JsonPath expression.
       
       When applied to an Array, the children to remove are those matching a second-level matching JsonPath (syntax following arrays conditional within the '[]' JsonPath construct)
       

  ```yaml
       # Syntax:
       
       - name_of_rule: whatever_delete_rule_name

         # match is a node selector, starting at your document root ($)
         
         #  Here the selector targets only some some 'containers' in my array
         # And for these containers, it will remove port 8080 from their ports list
       
         # For some hints on Jsonpath syntax, you can have a look for example at https://support.smartbear.com/alertsite/docs/monitors/api/endpoint/jsonpath.html

         match: "$..spec.containers[?(@.my_container_type_field==needs_ports_trimming)]" 
         deleteChildrenThatMatch: "?(@.ports[?(@.target_port==8080)])"

  ```

       When applied to a Map Object, this rule can also allow to remove a single key by its name. But in this case, wildcard or conditions based on the object content are not supported (only a single fixed key is removed)

  ```yaml
       # Syntax:
       
       - name_of_rule: whatever_delete_attribute_rule_name

         # match is a node selector, starting at your document root ($)
         
         #  Here the selector targets  all 'containers' items
         # and will remove their 'imagePolicy' attribute
         # And for these containers, it will remove port 8080 from their ports list
         match: "$..spec.containers[*]" 
         deleteChildrenThatMatch: "$.imagePolicy"

  ```  


As an example, see "rules_example.yaml", and "myfile_to_transform.yaml"


Usage:

       yaml-transform <rules yaml file path> <document yaml file path>

try out:
       
       ./yaml-transform  rules_example.yaml myfile_to_transform.yaml

This will output transformed yaml file to stdout

