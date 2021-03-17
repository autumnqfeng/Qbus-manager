package errno

// Error code description:
// 				1						00						02
// 		service level error		service module code		specific error code
//
// service level error: 1:system error, 2:common error
var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}

	// topic errors
	ErrCreateTopic       = &Errno{Code: 20101, Message: "The topic create failed."}
	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}

	// cluster errors
	ErrAddCluster            = &Errno{Code: 20201, Message: "The cluster add failed."}
	ErrClusterNotExist       = &Errno{Code: 20202, Message: "non-existing cluster"}
	ErrDisableCluster        = &Errno{Code: 20203, Message: "disable cluster error"}
	ErrDeleteCluster         = &Errno{Code: 20204, Message: "delete cluster error"}
	ErrNotFoundClusterConfig = &Errno{Code: 20204, Message: "not found cluster config"}
	ErrClusterConnect        = &Errno{Code: 20205, Message: "cluster connect failed."}
	ErrGetBroker             = &Errno{Code: 20206, Message: "get brokers from cluster failed."}

	// conn errors
	ErrAddConnPoll       = &Errno{Code: 20301, Message: "add cluster to connection pool failed."}
	ErrGetConn           = &Errno{Code: 20302, Message: "get cluster to connection failed."}
	ErrClusterMsgNotInKM = &Errno{Code: 20303, Message: "The cluster message not in kafka manager."}
)
