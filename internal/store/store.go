package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{ db *sql.DB }
type Flag struct { ID string `json:"id"`; Key string `json:"key"`; Name string `json:"name"`; Description string `json:"description,omitempty"`; Enabled bool `json:"enabled"`; Percentage int `json:"percentage"`; CreatedAt string `json:"created_at"`; UpdatedAt string `json:"updated_at"` }
func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil { return nil, err }
	db, err := sql.Open("sqlite", filepath.Join(d, "embargo.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil { return nil, err }
	db.Exec(`CREATE TABLE IF NOT EXISTS flags (id TEXT PRIMARY KEY, key TEXT UNIQUE NOT NULL, name TEXT DEFAULT '', description TEXT DEFAULT '', enabled INTEGER DEFAULT 0, percentage INTEGER DEFAULT 100, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')))`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) scan(s interface{Scan(...any)error}) *Flag {
	var f Flag; var en int; if s.Scan(&f.ID,&f.Key,&f.Name,&f.Description,&en,&f.Percentage,&f.CreatedAt,&f.UpdatedAt)!=nil{return nil}; f.Enabled=en==1; return &f
}
func (d *DB) Create(f *Flag) error { f.ID=genID();f.CreatedAt=now();f.UpdatedAt=f.CreatedAt; en:=0;if f.Enabled{en=1}; _,err:=d.db.Exec(`INSERT INTO flags VALUES(?,?,?,?,?,?,?,?)`,f.ID,f.Key,f.Name,f.Description,en,f.Percentage,f.CreatedAt,f.UpdatedAt); return err }
func (d *DB) Get(id string) *Flag { return d.scan(d.db.QueryRow(`SELECT * FROM flags WHERE id=?`,id)) }
func (d *DB) GetByKey(key string) *Flag { return d.scan(d.db.QueryRow(`SELECT * FROM flags WHERE key=?`,key)) }
func (d *DB) List() []Flag { rows,_:=d.db.Query(`SELECT * FROM flags ORDER BY key`); if rows==nil{return nil}; defer rows.Close(); var o []Flag; for rows.Next(){if f:=d.scan(rows);f!=nil{o=append(o,*f)}}; return o }
func (d *DB) Update(id string, f *Flag) error { en:=0;if f.Enabled{en=1}; _,err:=d.db.Exec(`UPDATE flags SET name=?,description=?,enabled=?,percentage=?,updated_at=? WHERE id=?`,f.Name,f.Description,en,f.Percentage,now(),id); return err }
func (d *DB) Toggle(id string) error { _,err:=d.db.Exec(`UPDATE flags SET enabled=1-enabled,updated_at=? WHERE id=?`,now(),id); return err }
func (d *DB) Delete(id string) error { _,err:=d.db.Exec(`DELETE FROM flags WHERE id=?`,id); return err }
func (d *DB) Evaluate(key string) (bool, *Flag) { f:=d.GetByKey(key); if f==nil||!f.Enabled{return false,f}; return true,f }
func (d *DB) Count() int { var n int; d.db.QueryRow(`SELECT COUNT(*) FROM flags`).Scan(&n); return n }
func (d *DB) EnabledCount() int { var n int; d.db.QueryRow(`SELECT COUNT(*) FROM flags WHERE enabled=1`).Scan(&n); return n }
