package fix

// StorageSide is a storage side.
type StorageSide string

// Storage sides for a message key.
const (
	Incoming StorageSide = "incoming"
	Outgoing StorageSide = "outgoing"
)

// StorageID is a storage id.
type StorageID struct {
	Sender string
	Target string
	Side   StorageSide
}
