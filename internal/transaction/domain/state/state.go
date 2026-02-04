package state

type State string

const (
    StateInitiated        State = "INITIATED"
    StateReserved         State = "RESERVED"
    StateOfflineCompleted State = "OFFLINE_COMPLETED"
    StatePendingSync      State = "PENDING_SYNC"
    StateValidating       State = "VALIDATING"
    StateCommitted        State = "COMMITTED"
    StateCompleted        State = "COMPLETED"
    StateRejected         State = "REJECTED"
    StateExpired          State = "EXPIRED"
    StateFailed           State = "FAILED"
)
