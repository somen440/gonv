/*
Copyright 2020 somen440

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package structure

import "strings"

// ViewStructure view structure
type ViewStructure struct {
	Name        string
	CreateQuery string
	Properties  []string
}

// CreateQueryToFormat format create query
func (vs *ViewStructure) CreateQueryToFormat() (query string) {
	query = vs.CreateQuery
	definer := vs.getDefiner()
	query = strings.Replace(query, definer, "", -1)
	query = strings.Replace(query, "AS select ", "\nAs select\n ", -1)
	query = strings.Replace(query, ",", ",\n ", -1)
	query = strings.Replace(query, " from ", " \nfrom\n ", -1)
	query = strings.Replace(query, " where ", " \nwhere\n ", -1)
	return
}

// Type return table structure type
func (vs *ViewStructure) Type() TableStructureType {
	return ViewRawType
}

// CompareQuery compare query
func (vs *ViewStructure) CompareQuery() (compareQuery string) {
	definer := vs.getDefiner()
	compareQuery = strings.Replace(vs.CreateQuery, definer, "", -1)
	compareQuery = strings.Replace(compareQuery, "\n", "", -1)
	compareQuery = strings.Replace(compareQuery, " ", "", -1)
	compareQuery = strings.Replace(compareQuery, vs.Name, "TABLENAME", -1)
	return
}

// GetName return Name
func (vs *ViewStructure) GetName() string {
	return vs.Name
}

func (vs *ViewStructure) getDefiner() string {
	return " DEFINER" + strings.Split(vs.CreateQuery, "DEFINER")[1] + "DEFINER"
}
