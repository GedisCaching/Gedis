package RESP

// WatchEventType is the type of event for watch events.

var (
	WatchSET = `
	    SET: is a function that sets a key-value pair in the store. 
		it takes a key and a value as arguments.
	    like this: SET key value
	    SET key value (EX, PX) 30 
	`

	WatchGET = `
	    GET: is a function that retrieves a value for a given key from the store.
	    it takes a key as argument.
	    like this: GET key
	`

	WatchDEL = `
	    DEL: is a function that deletes a key-value pair from the store.
	    it takes a key as argument.
	    like this: DEL key
	`

	WatchEXISTS = `
	    EXISTS: is a function that checks if a key exists in the store.
	    it takes a key as argument.
	    like this: EXISTS key
	    returns True if the key exists, False otherwise.
	`

	WatchTTL = `
	    TTL: is a function that returns the time to live of a key in seconds.
	    it takes a key as argument.
	    like this: TTL key
	    returns the remaining time in seconds, 0 if no expiry, or -2 if the key doesn't exist.
	`

	WatchPING = `
	    PING: is a function used to test if the server is responsive.
	    like this: PING
	    returns PONG
	`

	WatchEXPIRE = `
	    EXPIRE: is a function that sets a timeout on a key after which the key will be automatically deleted.
	    it takes a key and seconds as arguments.
	    like this: EXPIRE key seconds
	    returns 1 if the timeout was set, 0 if the key doesn't exist or the timeout couldn't be set.
	`

	WatchGETDEL = `
	    GETDEL: is a function that gets the value of a key and deletes the key.
	    it takes a key as argument.
	    like this: GETDEL key
	    returns the value of the key if it exists, or nil if the key doesn't exist.
	    this is an atomic operation that combines GET and DEL.
	`

	WatchRENAME = `
	    RENAME: is a function that renames a key to a new key name.
	    it takes the old key name and the new key name as arguments.
	    like this: RENAME oldkey newkey
	    returns OK if successful, or an error if the key doesn't exist or the new key already exists.
	    if newkey already exists, it is overwritten.
	`
)
