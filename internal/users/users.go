package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var dbMutex sync.Mutex

// UserAccess access
type UserAccess int

// UserAccess types
const (
	UANone UserAccess = iota
	UARevoked
	UAMember
	UAAdmin
)

// User struct object for access
type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Access    UserAccess `json:"access"`
}

// IsAdmin returns if a user is an admin or not
func (u User) IsAdmin() bool {
	return u.Access == UAAdmin
}

// IsMember returns if a user is a member or not
func (u User) IsMember() bool {
	return u.Access == UAMember
}

// IsRevoked returns if a user is revoked or not from memberlist
func (u User) IsRevoked() bool {
	return u.Access == UARevoked
}

// DisplayName returns username if set in telegram or first name + last name if not set
func (u User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// Recipient converts telegram id to string
func (u User) Recipient() string {
	return strconv.Itoa(u.ID)
}

// UserDB struct object for mapping users
type UserDB struct {
	users    []User
	usersMap map[int]User
	dbPath   string
}

// NewUserDB constructs user mapping struct from config
func NewUserDB(dbPath string) (db *UserDB, err error) {
	db = &UserDB{
		users:    []User{},
		usersMap: map[int]User{},
		dbPath:   dbPath,
	}

	loadErr := db.Load()
	if !os.IsNotExist(loadErr) {
		err = loadErr
	}

	return
}

// Create appends user info to user map
func (u *UserDB) Create(user User) error {
	if u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d already exists", user.ID)
	}

	u.users = append(u.users, user)
	u.usersMap[user.ID] = user

	u.Save()
	return nil
}

// Update updates already existing user info to user map
func (u *UserDB) Update(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d doesn't exist", user.ID)
	}

	for i := 0; i < len(u.users); i++ {
		if u.users[i].ID == user.ID {
			u.users[i] = user
			break
		}
	}

	u.usersMap[user.ID] = user
	u.Save()
	return nil
}

// Delete deletes already existing user info from user map
func (u *UserDB) Delete(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d doesn't exist", user.ID)
	}

	for i, usr := range u.users {
		if user.ID == usr.ID {
			u.users = append(u.users[:i], u.users[i+1:]...)
			break
		}
	}

	delete(u.usersMap, user.ID)
	u.Save()
	return nil
}

// User returns user id and whether it exists
func (u *UserDB) User(id int) (User, bool) {
	user, exists := u.usersMap[id]
	return user, exists
}

// Exists returns true if a user exists
func (u *UserDB) Exists(id int) bool {
	_, ok := u.usersMap[id]
	return ok
}

// Users returns all users in user map
func (u *UserDB) Users() []User {
	return u.users
}

// Admins returns list of admins in user map
func (u *UserDB) Admins() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAAdmin {
			result = append(result, user)
		}
	}
	return result
}

// Members returns all members in user map
func (u *UserDB) Members() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAMember {
			result = append(result, user)
		}
	}
	return result
}

// Revoked returns revoked users in user map
func (u *UserDB) Revoked() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UARevoked {
			result = append(result, user)
		}
	}
	return result
}

// IsAdmin returns whether a user is an admin or not
func (u *UserDB) IsAdmin(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAAdmin
}

// IsMember returns whether a user is a member or not
func (u *UserDB) IsMember(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAMember
}

// IsRevoked returns whether a user is revoked or not
func (u *UserDB) IsRevoked(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UARevoked
}

// Save saves changes to user map
func (u *UserDB) Save() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	var tempDb *os.File
	tempDb, err := ioutil.TempFile(filepath.Dir(u.dbPath), filepath.Base(u.dbPath))
	if err != nil {
		return err
	}

	var db = struct {
		Users []User `json:"users"`
	}{
		u.users,
	}

	enc := json.NewEncoder(tempDb)
	err = enc.Encode(db)
	if err != nil {
		return err
	}

	err = tempDb.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempDb.Name(), u.dbPath)
	if err != nil {
		return err
	}
	return nil
}

// Load loads user map from config file
func (u *UserDB) Load() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	raw, err := ioutil.ReadFile(u.dbPath)
	if err != nil {
		return err
	}

	var db struct {
		Users []User
	}

	json.Unmarshal(raw, &db)
	u.users = db.Users
	u.usersMap = map[int]User{}
	for _, user := range db.Users {
		u.usersMap[user.ID] = user
	}

	return nil
}
