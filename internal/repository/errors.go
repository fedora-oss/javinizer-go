package repository

import "errors"

// ErrNotFound is returned by repository methods when a requested record does
// not exist in the store.  Callers use errors.Is(err, ErrNotFound) to
// distinguish "missing record" from other database errors.
var ErrNotFound = errors.New("record not found")
