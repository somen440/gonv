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
	query = strings.Replace(query, definer, "", 0)
	query = strings.Replace(query, "AS select ", "\nAs select\n ", 0)
	query = strings.Replace(query, ",", ",\n ", 0)
	query = strings.Replace(query, " from ", " \nfrom\n ", 0)
	query = strings.Replace(query, " where ", " \nwhere\n ", 0)
	return
}

// Type return table structure type
func (vs *ViewStructure) Type() TableStructureType {
	return ViewRawType
}

// CompareQuery compare query
func (vs *ViewStructure) CompareQuery() (compareQuery string) {
	definer := vs.getDefiner()
	compareQuery = strings.Replace(vs.CreateQuery, definer, "", 0)
	compareQuery = strings.Replace(compareQuery, "\n", "", 0)
	compareQuery = strings.Replace(compareQuery, " ", "", 0)
	compareQuery = strings.Replace(compareQuery, vs.Name, "TABLENAME", 0)
	return
}

// GetName return Name
func (vs *ViewStructure) GetName() string {
	return vs.Name
}

func (vs *ViewStructure) getDefiner() string {
	return " DEFINER" + strings.Split(vs.CreateQuery, "DEFINER")[1] + "DEFINER"
}
