package context

import (
	// "fmt"
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
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	Uid        string    `json:"uid"`
	Generated  bool      `json:"generated"`
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

	if err := r.Scan(&URL, &CreatedAt, &ModifiedAt, &DeletedAt, &Uid, &Generated, &Name); err != nil {
		/*Scan's destinations have to be in the same order as the columns in the schema*/
		return nil, "", err
	}

	rt := &Route{URL: URL, CreatedAt: CreatedAt, ModifiedAt: ModifiedAt, DeletedAt: DeletedAt, Uid: Uid, Generated: Generated}

	return rt, Name, nil
}

func createTableIfNotExist(db *sql.DB) error {
	// if a table called linkdata does not exist, set it up
	queryString := "CREATE TABLE IF NOT EXISTS linkdata (URL varchar(500) NOT NULL, CreatedAt timestamp NOT NULL, ModifiedAt timestamp, DeletedAt timestamp, Uid bigint PRIMARY KEY, Generated boolean NOT NULL, Name varchar(100) NOT NULL)"
	_, err := db.Exec(queryString)

	return err
}

func dropTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE linkdata")
	return err
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
		err = dropTable(db)
	}

	err = createTableIfNotExist(db)
	if err != nil {
		return nil, err
	}

	return &Context{
		db: db,
	}, nil
}

// Context provides access to the database.
type Context struct {
	db *sql.DB
}

// Close the database associated with this context.
func (c *Context) Close() error {
	return c.db.Close()
}

// GetUid retreives a single shortcut matching 'id' from the data store.
func (c *Context) GetUid(uid string) (*Route, error) {
	// the row returned from the database should have the same number of fields (with the same names) as the fields in the definition of the Route object.
	rows, err := c.db.Query("SELECT * FROM linkdata WHERE Uid = $1", uid)
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
	rows, err := c.db.Query("SELECT * FROM linkdata WHERE Name = $1", name)
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
func (c *Context) Edit(uid, newName, newUrl string) error {
	_, err := c.db.Exec("UPDATE linkdata SET Url = $1, ModifiedAt = $2, Name = $3 WHERE Uid = $4", newUrl, time.Now(), newName, uid)

	return err
}

// Put creates a new row from a route and a name and inserts it into the database.
/*
What if someone wants to edit an existing row by name?
Should we check that in api/apiUrlPost (which currently generates a new Route and calls ctx.Put) or here?
Probably we should have a different method here like Edit, just so these are easy to handle.
*/
func (c *Context) Put(name string, rt *Route) error {
	_, err := c.db.Exec("INSERT INTO linkdata VALUES ($1, $2, $3, $4, $5, $6, $7)", rt.URL, rt.CreatedAt, rt.ModifiedAt, rt.DeletedAt, rt.Uid, rt.Generated, name)

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

	// var URL string
	// var Time time.Time
	// var Uid uint32
	// var Generated bool
	// var Name string

	for rows.Next() {
		rt, rowName, err := rowToRoute(rows)

		if err != nil {
			return nil, err
		}

		golinks[rowName] = *rt

	}

	return golinks, nil
}
