package iaaseki

const controllerAgentName = "EkiMonitor"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a EkiMonitor fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by EkiMonitor"
	// MessageResourceSynced is the message used for an Event fired when a EkiMonitor
	// is synced successfully
	MessageResourceSynced = "EkiMonitor synced successfully"
)
