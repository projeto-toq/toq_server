package workflow

import "errors"

// ErrFinalizationAccessDenied indicates the backend IAM role cannot start the
// Step Functions execution responsible for media finalization.
var ErrFinalizationAccessDenied = errors.New("workflow finalization access denied")
