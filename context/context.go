package context

import (
	"fmt"
	"math/rand"
	"strings"

	"os"
	"time"

	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

/*NOTES: List method removed because refactoring /api/urls/ seems to be the most work for least reward*/

const TABLE_NAME = "linkdata"

// Route is the value part of a shortcut.
type Route struct {
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
	DeletedAt     time.Time `json:"deleted_at"`
	Uid           string    `json:"uid"`
	Generated     bool      `json:"generated"`
	ModifiedCount int       `json:"modified_count"`
	/*A field declaration may be followed by an optional string literal tag, which becomes an attribute for all the fields in the corresponding field declaration.
	  The tags are made visible through a reflection interface and take part in type identity for structs but are otherwise ignored...*/
}

/*Refactor to return a name as well, probably. We may need a different version that takes sql.Rows instead of sql.Row .*/

// Takes a Rows object returned from a database query, repackages it into a Route, and returns that plus a name.
// for rows.Next() { rt = rowToRoute(rows) } should be the proper usage, I think.
func rowToRoute(r *sql.Rows) (*Route, string, error) {
	var URL string
	var CreatedAt time.Time
	var ModifiedAt time.Time
	var DeletedAt time.Time
	var Uid string
	var Generated bool
	var Name string
	var ModifiedCount int

	if err := r.Scan(&URL, &CreatedAt, &ModifiedAt, &DeletedAt, &Uid, &Generated, &Name, &ModifiedCount); err != nil {
		/*Scan's destinations have to be in the same order as the columns in the schema*/
		return nil, "", err
	}

	rt := &Route{URL: URL, CreatedAt: CreatedAt, ModifiedAt: ModifiedAt, DeletedAt: DeletedAt, Uid: Uid, Generated: Generated, ModifiedCount: ModifiedCount}

	return rt, Name, nil
}

func createTableIfNotExist(db *sql.DB, name string) error {
	// if a table called name does not exist, set it up
	queryString := "CREATE TABLE IF NOT EXISTS " + name + " (URL varchar(500) NOT NULL, CreatedAt timestamp NOT NULL, ModifiedAt timestamp, DeletedAt timestamp, Uid varchar(100) PRIMARY KEY, Generated boolean NOT NULL, Name varchar(100) NOT NULL, ModifiedCount int NOT NULL)"
	_, err := db.Exec(queryString)

	return err
}

func dropTable(db *sql.DB, name string) error {
	_, err := db.Exec("DROP TABLE " + name)
	return err
}

func (c *Context) DropTable() error {
	if !strings.Contains(c.table_name, "test") {
		return errors.New("This context does not appear to be a test context!")
	} else {
		return dropTable(c.db, c.table_name)
	}
}

// Creates a Context that contains a sql.DB (postgres database) and returns a pointer to said context.
func Open() (*Context, error) {
	// open the database and return db, a pointer to the sql.DB object
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	// ping the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if os.Getenv("DROPTABLE_EACH_LAUNCH") == "yes" { /*Turn this off once we're ready to launch*/
		err = dropTable(db, TABLE_NAME)
	}

	err = createTableIfNotExist(db, TABLE_NAME)
	if err != nil {
		return nil, err
	}

	return &Context{
		table_name: TABLE_NAME,
		db:         db,
	}, nil
}

func OpenTestCtx() (*Context, error) {
	// open the database and return db, a pointer to the sql.DB object
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}

	// ping the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	table_name := fmt.Sprintf("olinks_test_%v", rand.Uint64())

	err = createTableIfNotExist(db, table_name)
	if err != nil {
		return nil, err
	}

	return &Context{
		db:         db,
		table_name: table_name,
	}, nil
}

// Context provides access to the database.
type Context struct {
	db         *sql.DB
	table_name string
}

// Close the database associated with this context.
func (c *Context) Close() error {
	return c.db.Close()
}

// GetUid retreives a single shortcut matching 'id' from the data store.
func (c *Context) GetUid(uid string) (*Route, error) {
	// the row returned from the database should have the same number of fields (with the same names) as the fields in the definition of the Route object.
	rows, err := c.db.Query("SELECT * FROM "+c.table_name+" WHERE Uid = $1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rt, _, err := rowToRoute(rows)
		if err != nil {
			return nil, err
		}
		return rt, nil
	}
	return nil, sql.ErrNoRows
}

// Get retreives a single shortcut matching 'name' from the data store.
func (c *Context) Get(name string) (*Route, error) {
	// the row returned from the database should have the same number of fields (with the same names) as the fields in the definition of the Route object.
	rows, err := c.db.Query("SELECT * FROM "+c.table_name+" WHERE Name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rt, _, err := rowToRoute(rows)
		if err != nil {
			return nil, err
		}
		return rt, nil
	}
	return nil, sql.ErrNoRows
}

//Edits the name and URL of a row and updates the ModifiedAt timestamp accordingly. Might want to generalize in the future.
func (c *Context) Edit(route *Route, name string) error {
	_, err := c.db.Exec("UPDATE "+c.table_name+" SET Url = $1, ModifiedAt = $2, Name = $3, ModifiedCount = $4, Generated=$5 WHERE Uid = $6",
		route.URL, time.Now().In(time.UTC), name, route.ModifiedCount, route.Generated, route.Uid)

	return err
}

// Put creates a new row from a route and a name and inserts it into the database.
/*
What if someone wants to edit an existing row by name?
Should we check that in api/apiUrlPost (which currently generates a new Route and calls ctx.Put) or here?
Probably we should have a different method here like Edit, just so these are easy to handle.
*/
func (c *Context) Put(name string, rt *Route) error {
	_, err := c.db.Exec("INSERT INTO "+c.table_name+" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		rt.URL, rt.CreatedAt.In(time.UTC), rt.ModifiedAt.In(time.UTC), rt.DeletedAt.In(time.UTC), rt.Uid, rt.Generated, name, rt.ModifiedCount)

	return err
}

// Del removes an existing shortcut from the data store.
func (c *Context) Del(name string) error {

	_, err := c.db.Exec("DELETE FROM "+c.table_name+" WHERE Name = $1", name)

	return err
}

// GetAll gets everything in the db to dump it out for backup purposes
func (c *Context) GetAll() (map[string]Route, error) {
	golinks := map[string]Route{}

	rows, err := c.db.Query("SELECT * FROM " + c.table_name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		rt, rowName, err := rowToRoute(rows)

		if err != nil {
			return nil, err
		}

		golinks[rowName] = *rt

	}

	return golinks, nil
}
