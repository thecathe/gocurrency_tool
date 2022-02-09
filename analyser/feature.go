package main

type FeatureType string

const (
	NONE                                              FeatureType = "None"
	NONE_COUNT                                        int         = 0
	GOROUTINE                                         FeatureType = "Goroutine"
	GOROUTINE_COUNT                                   int         = 1
	RECEIVE                                           FeatureType = "Receive"
	RECEIVE_COUNT                                     int         = 2
	SEND                                              FeatureType = "Send"
	SEND_COUNT                                        int         = 3
	MAKE_CHAN                                         FeatureType = "Synchronous chan"
	MAKE_CHAN_COUNT                                   int         = 4
	GO_IN_FOR                                         FeatureType = "Go in for"
	GO_IN_FOR_COUNT                                   int         = 5
	RANGE_OVER_CHAN                                   FeatureType = "Range over chan"
	RANGE_OVER_CHAN_COUNT                             int         = 6
	GO_IN_CONSTANT_FOR                                FeatureType = "Goroutine in for with constant (constant)"
	GO_IN_CONSTANT_FOR_COUNT                          int         = 7
	KNOWN_CHAN_DEPTH                                  FeatureType = "Known chan length (length)"
	KNOWN_CHAN_DEPTH_COUNT                            int         = 8
	UNKNOWN_CHAN_DEPTH                                FeatureType = "Unknown chan length"
	UNKNOWN_CHAN_DEPTH_COUNT                          int         = 9
	MAKE_CHAN_IN_FOR                                  FeatureType = "Make chan in for"
	MAKE_CHAN_IN_FOR_COUNT                            int         = 10
	MAKE_CHAN_IN_CONSTANT_FOR                         FeatureType = "Make chan in constant for"
	MAKE_CHAN_IN_CONSTANT_FOR_COUNT                   int         = 23
	ARRAY_OF_CHANNELS                                 FeatureType = "Array of chans"
	ARRAY_OF_CHANNELS_COUNT                           int         = 11
	CONSTANT_CHAN_ARRAY                               FeatureType = "Constant array of chans (length)"
	CONSTANT_CHAN_ARRAY_COUNT                         int         = 12
	CHAN_SLICE                                        FeatureType = "Slice array of chans"
	CHAN_SLICE_COUNT                                  int         = 13
	CHAN_MAP                                          FeatureType = "Map of chans"
	CHAN_MAP_COUNT                                    int         = 14
	CLOSE_CHAN                                        FeatureType = "Close chan"
	CLOSE_CHAN_COUNT                                  int         = 15
	SELECT                                            FeatureType = "Select (number of branch)"
	SELECT_COUNT                                      int         = 16
	DEFAULT_SELECT                                    FeatureType = "Select with default (number of branch)"
	DEFAULT_SELECT_COUNT                              int         = 17
	ASSIGN_CHAN_IN_FOR                                FeatureType = "Assign chan in for"
	ASSIGN_CHAN_IN_FOR_COUNT                          int         = 18
	CHAN_OF_CHANS                                     FeatureType = "Channel of channels"
	CHAN_OF_CHANS_COUNT                               int         = 19
	RECEIVE_CHAN                                      FeatureType = "Receive only chan (<-chan)"
	RECEIVE_CHAN_COUNT                                int         = 20
	SEND_CHAN                                         FeatureType = "Send only chan (chan<-)"
	SEND_CHAN_COUNT                                   int         = 21
	PARAM_CHAN                                        FeatureType = "chan used as a param"
	PARAM_CHAN_COUNT                                  int         = 22
	WAITGROUP                                         FeatureType = "Waitgroup"
	WAITGROUP_COUNT                                   int         = 23
	KNOWN_ADD                                         FeatureType = "Waitgroup Add(const)"
	KNOWN_ADD_COUNT                                   int         = 24
	UNKNOWN_ADD                                       FeatureType = "Waitgroup Add(var)"
	UNKNOWN_ADD_COUNT                                 int         = 25
	DONE                                              FeatureType = "Waitgroup Done()"
	DONE_COUNT                                        int         = 26
	MUTEX                                             FeatureType = "Mutex"
	MUTEX_COUNT                                       int         = 27
	UNLOCK                                            FeatureType = "Mutex Unlock()"
	UNLOCK_COUNT                                      int         = 28
	LOCK                                              FeatureType = "Mutex Lock()"
	LOCK_COUNT                                        int         = 29
	SELECT_SYNC_R_EXCL_TIMEOUT                        FeatureType = "Select with only synchronous receive, excluding timeouts (number of branch)"
	SELECT_SYNC_R_EXCL_TIMEOUT_COUNT                  int         = 30
	SELECT_SYNC_R_INCL_TIMEOUT                        FeatureType = "Select with only synchronous receive, including timeouts (number of branch)"
	SELECT_SYNC_R_INCL_TIMEOUT_COUNT                  int         = 31
	SELECT_SYNC_R_DEFAULT_EXCL_TIMEOUT                FeatureType = "Select with only synchronous receive, excluding timeouts, and default (number of branch)"
	SELECT_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT          int         = 32
	SELECT_SYNC_R_DEFAULT_INCL_TIMEOUT                FeatureType = "Select with only synchronous receive, including timeouts, and default (number of branch)"
	SELECT_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT          int         = 33
	SELECT_ASYNC_R_EXCL_TIMEOUT                       FeatureType = "Select with only aynchronous receive, excluding timeouts (number of branch)"
	SELECT_ASYNC_R_EXCL_TIMEOUT_COUNT                 int         = 34
	SELECT_ASYNC_R_INCL_TIMEOUT                       FeatureType = "Select with only aynchronous receive, including timeouts (number of branch)"
	SELECT_ASYNC_R_INCL_TIMEOUT_COUNT                 int         = 35
	SELECT_ASYNC_R_DEFAULT_EXCL_TIMEOUT               FeatureType = "Select with only aynchronous receive, excluding timeouts, and default (number of branch)"
	SELECT_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT         int         = 36
	SELECT_ASYNC_R_DEFAULT_INCL_TIMEOUT               FeatureType = "Select with only aynchronous receive, including timeouts, and default (number of branch)"
	SELECT_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT         int         = 37
	SELECT_SYNC_S                                     FeatureType = "Select with only synchronous send (number of branches)"
	SELECT_SYNC_S_COUNT                               int         = 38
	SELECT_SYNC_S_DEFAULT                             FeatureType = "Select with only synchronous send and default (number of branches)"
	SELECT_SYNC_S_DEFAULT_COUNT                       int         = 39
	SELECT_SYNC_S_TIMEOUT                             FeatureType = "Select with only synchronous send and timeouts (number of branches)"
	SELECT_SYNC_S_TIMEOUT_COUNT                       int         = 40
	SELECT_SYNC_S_DEFAULT_TIMEOUT                     FeatureType = "Select with only synchronous send, default and timeouts (number of branches)"
	SELECT_SYNC_S_DEFAULT_TIMEOUT_COUNT               int         = 41
	SELECT_ASYNC_S                                    FeatureType = "Select with only asynchronous send (number of branches)"
	SELECT_ASYNC_S_COUNT                              int         = 42
	SELECT_ASYNC_S_DEFAULT                            FeatureType = "Select with only asynchronous send and default (number of branches)"
	SELECT_ASYNC_S_DEFAULT_COUNT                      int         = 43
	SELECT_ASYNC_S_TIMEOUT                            FeatureType = "Select with only asynchronous send and timeouts (number of branches)"
	SELECT_ASYNC_S_TIMEOUT_COUNT                      int         = 44
	SELECT_ASYNC_S_DEFAULT_TIMEOUT                    FeatureType = "Select with only asynchronous send, default and timeouts (number of branches)"
	SELECT_ASYNC_S_DEFAULT_TIMEOUT_COUNT              int         = 45
	SELECT_SYNC_S_SYNC_R_EXCL_TIMEOUT                 FeatureType = "Select with only synchronous send and synchronous receive, excluding timeouts (number of branches)"
	SELECT_SYNC_S_SYNC_R_EXCL_TIMEOUT_COUNT           int         = 46
	SELECT_SYNC_S_SYNC_R_INCL_TIMEOUT                 FeatureType = "Select with only synchronous send and synchronous receive, including timeouts (number of branches)"
	SELECT_SYNC_S_SYNC_R_INCL_TIMEOUT_COUNT           int         = 47
	SELECT_SYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT         FeatureType = "Select with only synchronous send, synchronous receive and defaults, excluding timeouts (number of branches)"
	SELECT_SYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT   int         = 48
	SELECT_SYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT         FeatureType = "Select with only synchronous send, synchronous receive and defaults, including timeouts (number of branches)"
	SELECT_SYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT   int         = 49
	SELECT_ASYNC_S_SYNC_R_EXCL_TIMEOUT                FeatureType = "Select with only asynchronous send and synchronous receive, excluding timeouts (number of branches)"
	SELECT_ASYNC_S_SYNC_R_EXCL_TIMEOUT_COUNT          int         = 50
	SELECT_ASYNC_S_SYNC_R_INCL_TIMEOUT                FeatureType = "Select with only asynchronous send and synchronous receive, including timeouts (number of branches)"
	SELECT_ASYNC_S_SYNC_R_INCL_TIMEOUT_COUNT          int         = 51
	SELECT_ASYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT        FeatureType = "Select with only asynchronous send, synchronous receive and defaults, excluding timeouts (number of branches)"
	SELECT_ASYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT  int         = 52
	SELECT_ASYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT        FeatureType = "Select with only asynchronous send, synchronous receive and defaults, including timeouts (number of branches)"
	SELECT_ASYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT  int         = 53
	SELECT_SYNC_S_ASYNC_R_EXCL_TIMEOUT                FeatureType = "Select with only synchronous send and asynchronous receive, excluding timeouts (number of branches)"
	SELECT_SYNC_S_ASYNC_R_EXCL_TIMEOUT_COUNT          int         = 54
	SELECT_SYNC_S_ASYNC_R_INCL_TIMEOUT                FeatureType = "Select with only synchronous send and asynchronous receive, including timeouts (number of branches)"
	SELECT_SYNC_S_ASYNC_R_INCL_TIMEOUT_COUNT          int         = 55
	SELECT_SYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT        FeatureType = "Select with only synchronous send, asynchronous receive and defaults, excluding timeouts (number of branches)"
	SELECT_SYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT  int         = 56
	SELECT_SYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT        FeatureType = "Select with only synchronous send, asynchronous receive and defaults, including timeouts (number of branches)"
	SELECT_SYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT  int         = 57
	SELECT_ASYNC_S_ASYNC_R_EXCL_TIMEOUT               FeatureType = "Select with only asynchronous send and asynchronous receive, excluding timeouts (number of branches)"
	SELECT_ASYNC_S_ASYNC_R_EXCL_TIMEOUT_COUNT         int         = 58
	SELECT_ASYNC_S_ASYNC_R_INCL_TIMEOUT               FeatureType = "Select with only asynchronous send and asynchronous receive, including timeouts (number of branches)"
	SELECT_ASYNC_S_ASYNC_R_INCL_TIMEOUT_COUNT         int         = 59
	SELECT_ASYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT       FeatureType = "Select with only asynchronous send, asynchronous receive and defaults, excluding timeouts (number of branches)"
	SELECT_ASYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT int         = 60
	SELECT_ASYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT       FeatureType = "Select with only asynchronous send, asynchronous receive and defaults, including timeouts (number of branches)"
	SELECT_ASYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT int         = 61
)

type Feature struct {
	F_type         FeatureType
	F_type_num     int
	F_filename     string
	F_package_name string
	F_line_num     int
	F_number       string // A number used to report additional info about a feature
	F_commit       string // the commit of the project at the time the feature was found
	F_project_name string // the project name of the feature
}

// takes a list of feature and sets their feature number according to their types
func setFeaturesNumber(counter *Counter) {

	features := counter.Features
	counter.Features = []*Feature{}

	for _, feature := range features {
		switch feature.F_type {
		case NONE:
			feature.F_type_num = NONE_COUNT
		case GOROUTINE:
			feature.F_type_num = GOROUTINE_COUNT
		case RECEIVE:
			feature.F_type_num = RECEIVE_COUNT
		case SEND:
			feature.F_type_num = SEND_COUNT
		case MAKE_CHAN:
			feature.F_type_num = MAKE_CHAN_COUNT
		case GO_IN_FOR:
			feature.F_type_num = GO_IN_FOR_COUNT
		case RANGE_OVER_CHAN:
			feature.F_type_num = RANGE_OVER_CHAN_COUNT
		case GO_IN_CONSTANT_FOR:
			feature.F_type_num = GO_IN_CONSTANT_FOR_COUNT
		case KNOWN_CHAN_DEPTH:
			feature.F_type_num = KNOWN_CHAN_DEPTH_COUNT
		case UNKNOWN_CHAN_DEPTH:
			feature.F_type_num = UNKNOWN_CHAN_DEPTH_COUNT
		case MAKE_CHAN_IN_FOR:
			feature.F_type_num = MAKE_CHAN_IN_FOR_COUNT
		case MAKE_CHAN_IN_CONSTANT_FOR:
			feature.F_type_num = MAKE_CHAN_IN_CONSTANT_FOR_COUNT
		case ARRAY_OF_CHANNELS:
			feature.F_type_num = ARRAY_OF_CHANNELS_COUNT
		case CONSTANT_CHAN_ARRAY:
			feature.F_type_num = CONSTANT_CHAN_ARRAY_COUNT
		case CHAN_SLICE:
			feature.F_type_num = CHAN_SLICE_COUNT
		case CHAN_MAP:
			feature.F_type_num = CHAN_MAP_COUNT
		case CLOSE_CHAN:
			feature.F_type_num = CLOSE_CHAN_COUNT
		case SELECT:
			feature.F_type_num = SELECT_COUNT
		case DEFAULT_SELECT:
			feature.F_type_num = DEFAULT_SELECT_COUNT
		case ASSIGN_CHAN_IN_FOR:
			feature.F_type_num = ASSIGN_CHAN_IN_FOR_COUNT
		case CHAN_OF_CHANS:
			feature.F_type_num = CHAN_OF_CHANS_COUNT
		case RECEIVE_CHAN:
			feature.F_type_num = RECEIVE_CHAN_COUNT
		case SEND_CHAN:
			feature.F_type_num = SEND_CHAN_COUNT
		case PARAM_CHAN:
			feature.F_type_num = PARAM_CHAN_COUNT
		case WAITGROUP:
			feature.F_type_num = WAITGROUP_COUNT
		case KNOWN_ADD:
			feature.F_type_num = KNOWN_ADD_COUNT
		case UNKNOWN_ADD:
			feature.F_type_num = UNKNOWN_ADD_COUNT
		case DONE:
			feature.F_type_num = DONE_COUNT
		case MUTEX:
			feature.F_type_num = MUTEX_COUNT
		case UNLOCK:
			feature.F_type_num = UNLOCK_COUNT
		case LOCK:
			feature.F_type_num = LOCK_COUNT
		case SELECT_SYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_SYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		case SELECT_SYNC_S:
			feature.F_type_num = SELECT_SYNC_S_COUNT
		case SELECT_SYNC_S_DEFAULT:
			feature.F_type_num = SELECT_SYNC_S_DEFAULT_COUNT
		case SELECT_SYNC_S_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_TIMEOUT_COUNT
		case SELECT_SYNC_S_DEFAULT_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_DEFAULT_TIMEOUT_COUNT
		case SELECT_ASYNC_S:
			feature.F_type_num = SELECT_ASYNC_S_COUNT
		case SELECT_ASYNC_S_DEFAULT:
			feature.F_type_num = SELECT_ASYNC_S_DEFAULT_COUNT
		case SELECT_ASYNC_S_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_TIMEOUT_COUNT
		case SELECT_ASYNC_S_DEFAULT_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_DEFAULT_TIMEOUT_COUNT
		case SELECT_SYNC_S_SYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_SYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_SYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_SYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_SYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_SYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_SYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_SYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_SYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_SYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_ASYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_ASYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_ASYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_ASYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_SYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_SYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_ASYNC_R_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_ASYNC_R_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_ASYNC_R_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_ASYNC_R_INCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_ASYNC_R_DEFAULT_EXCL_TIMEOUT_COUNT
		case SELECT_ASYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT:
			feature.F_type_num = SELECT_ASYNC_S_ASYNC_R_DEFAULT_INCL_TIMEOUT_COUNT
		}

		counter.Features = append(counter.Features, feature)
	}
}
