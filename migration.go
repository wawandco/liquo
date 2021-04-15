package liquo

// Migration xml with liquibase format. A migration may be composed
// of multiple changesets.
type Migration struct {
	ChangeSets []ChangeSet `xml:"changeSet"`
}
