package liquo_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wawandco/liquo"
)

func TestReadMigration(t *testing.T) {
	r := require.New(t)
	data := `
		<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:ext="http://www.liquibase.org/xml/ns/dbchangelog-ext" xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-3.0.xsd http://www.liquibase.org/xml/ns/dbchangelog-ext http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-ext.xsd">
		<changeSet id="20210203002030-create_org_units" author="ox">
			<sql>	 
				CREATE TABLE organizational_units (
				id 		 uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
				name 		 character varying(255) NOT NULL,
				address 	 character varying(255) DEFAULT '' NOT NULL,
				active 	 boolean DEFAULT false NOT NULL,
				parent_id  uuid NULL,
				created_at timestamp without time zone NOT NULL,
				updated_at timestamp without time zone NOT NULL
				);

				COMMENT ON TABLE organizational_units IS 'This table holds the organizational units and their relationships. From companies to small departments live here.';
				COMMENT ON COLUMN organizational_units.name IS 'The name of the org unit.';
				COMMENT ON COLUMN organizational_units.active IS 'Define if its active or not';
				COMMENT ON COLUMN organizational_units.address IS 'Address in case of being a company';
				COMMENT ON COLUMN organizational_units.parent_id IS 'The id of the parent org unit';
				COMMENT ON COLUMN organizational_units.created_at IS 'timestamp when created';
				COMMENT ON COLUMN organizational_units.updated_At IS 'timestamp when updated';
			</sql>
			<sql>SELECT 1;</sql>
			<rollback>
				DROP TABLE IF EXISTS organizational_units;
			</rollback>
		</changeSet>
	</databaseChangeLog>
		`

	dir := t.TempDir()
	filename := filepath.Join(dir, "migration.xml")
	err := ioutil.WriteFile(filename, []byte(data), 0777)
	r.NoError(err, "could not create file")

	c := &liquo.Command{}
	m, err := c.ReadMigration(filename)
	r.NoError(err, "could not parse migration")
	r.NotNil(m)

	r.Len(m.ChangeSets, 1)
	r.Len(m.ChangeSets[0].SQL, 2)
	r.Contains(m.ChangeSets[0].SQL[0], `CREATE TABLE organizational_units (`)
	r.Contains(m.ChangeSets[0].SQL[1], `SELECT 1`)
}
