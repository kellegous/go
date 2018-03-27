package context

import (
	"os"
	"time"

	"database/sql"
	_ "github.com/lib/pq"
)

const (
	routesDbFilename = "routes.db"
)

/*NOTES: List method removed because refactoring /api/urls/ seems to be the most work for least reward*/

// Route is the value part of a shortcut.
type Route struct {
	URL       string    `json:"url"`
	Time      time.Time `json:"time"`
	Uid       uint32    `json:"uid"`
	Generated bool      `json:"generated"`
	/*A field declaration may be followed by an optional string literal tag, which becomes an attribute for all the fields in the corresponding field declaration.
	  The tags are made visible through a reflection interface and take part in type identity for structs but are otherwise ignored...*/
}

// Takes a Row object returned from a database query and repackages it into a Route.
func rowToRoute(r *sql.Row) (*Route, error) {
	var URL string
	var Time time.Time
	var Uid uint32
	var Generated bool
	var Name string

	if err := r.Scan(&URL, &Time, &Uid, &Generated, &Name); err != nil {
		/*Scan's destinations have to be in the same order as the columns in the schema*/
		return nil, err
	}

	rt := &Route{URL, Time, Uid, Generated}

	return rt, nil
}

func createTableIfNotExist(db *sql.DB) error {
	// if a table called linkdata does not exist, set it up
	queryString := "CREATE TABLE IF NOT EXISTS linkdata (URL varchar(500) NOT NULL, Time date NOT NULL, Uid bigint PRIMARY KEY, Generated boolean NOT NULL, Name varchar(100) NOT NULL)"
	_, err := db.Exec(queryString)

	return err
}

func dropTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE linkdata")
	return err
}

// Creates a Context that contains a sql.DB (postgres database) and returns a pointer to said context.
// Currently path isn't used for anything.
func Open(path string) (*Context, error) {
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

	err = createTableIfNotExist(db)
	if err != nil {
		return nil, err
	}

	return &Context{
		path: path,
		db:   db,
	}, nil
}

// Context provides access to the database. path is unnecessary now.
type Context struct {
	path string
	db   *sql.DB // possibly just *DB, I'm not 100% sure here
}

// Close the database associated with this context.
func (c *Context) Close() error {
	return c.db.Close()
}

// Get retreives a shortcut matching 'name' from the data store.
func (c *Context) Get(name string) (*Route, error) {

	// the row returned from the database should have the same number of fields (with the same names) as the fields in the definition of the Route object.
	row := c.db.QueryRow("SELECT * FROM linkdata WHERE Name = $1", name)

	return rowToRoute(row)
}

// Put creates a new row from a route and a name and inserts it into the database.
func (c *Context) Put(name string, rt *Route) error {
	_, err := c.db.Exec("INSERT INTO linkdata VALUES ($1, $2, $3, $4, $5)", rt.URL, rt.Time, rt.Uid, rt.Generated, name)

	return err
}

// Del removes an existing shortcut from the data store.
func (c *Context) Del(name string) error {

	_, err := c.db.Exec("DELETE FROM linkdata WHERE Name = $1", name)

	return err
}

// GetAll gets everything in the db to dump it out for backup purposes
func (c *Context) GetAll() (map[string]Route, error) {
	golinks := map[string]Route{}

	rows, err := c.db.Query("SELECT * FROM linkdata")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var URL string
	var Time time.Time
	var Uid uint32
	var Generated bool
	var Name string

	for rows.Next() {

		if err := rows.Scan(&URL, &Time, &Uid, &Generated, &Name); err != nil {
			return nil, err
		}

		rt := &Route{URL, Time, Uid, Generated}
		golinks[Name] = *rt
	}

	return golinks, nil
}
