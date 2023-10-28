package cli

// type Backend string

// const (
// 	BackendLevelDB   Backend = "leveldb"
// 	BackendFirestore Backend = "firestore"
// )

// func (b Backend) String() string {
// 	return string(b)
// }

// func (b *Backend) Set(v string) error {
// 	be, err := backendFromString(v)
// 	if err != nil {
// 		return err
// 	}
// 	*b = be
// 	return nil
// }

// func (b Backend) Type() string {
// 	return "string"
// }

// func backendFromString(s string) (Backend, error) {
// 	switch s {
// 	case BackendFirestore.String():
// 		return BackendFirestore, nil
// 	case BackendLevelDB.String():
// 		return BackendLevelDB, nil
// 	}
// 	return "", fmt.Errorf("invalid backend: %s", s)
