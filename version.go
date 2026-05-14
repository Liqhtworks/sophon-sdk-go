package sophon

// Version is the SDK release version. Kept in sync with the top-level
// VERSION file by the publish pipeline; baked into the default User-Agent
// so server-side logs can identify the SDK version that made each call.
const Version = "0.1.5"

// UserAgent is the recommended default User-Agent value. Set
//
//	cfg.UserAgent = sophon.UserAgent
//
// when constructing a Configuration to identify your SDK build in logs.
const UserAgent = "sophon-sdk-go/" + Version
