package app

const (
	STATE_FOCUS_URL             = "state_focus_url"
	STATE_FOCUS_RESPONSE        = "state_focus_response"
	STATE_FOCUS_PIPE            = "state_focus_pipe"
	STATE_FOCUS_PIPEDRESP       = "state_focus_pipedresp"
	STATE_FOCUS_RESPONSE_FILTER = "state_focus_response_filter"
	STATE_FOCUS_BODY            = "state_focus_body"
	STATE_FOCUS_HEADER          = "state_focus_header"
	STATE_COMMAND_PALLETE       = "state_command_pallete"
	STATE_METHOD_PALLETE        = "state_method_pallete"
	STATE_SELECT_ENV            = "state_select_env"
	STATE_SELECT_SESSION        = "state_select_session"
	STATE_SAVE_SESSION          = "state_save_session"
	STATE_SAVE_SESSION_INPUT    = "state_save_session_input"
	STATE_SESSION_RENAME_INPUT  = "state_session_rename_input"

	COMMAND_OPEN_ENV          = "command_open_env"
	COMMAND_CHANGE_ENV        = "command_change_env"
	COMMAND_OPEN_BODY         = "command_open_body"
	COMMAND_OPEN_HEADER       = "command_open_header"
	COMMAND_SELECT_METHOD     = "command_select_method"
	COMMAND_SAVE_SESSION      = "command_save_session"
	COMMAND_OPEN_SESSION_LIST = "command_open_session_list"

	REQUEST_METHOD_GET     = "GET"
	REQUEST_METHOD_POST    = "POST"
	REQUEST_METHOD_PUT     = "PUT"
	REQUEST_METHOD_DELETE  = "DELETE"
	REQUEST_METHOD_OPTIONS = "OPTIONS"
	REQUEST_METHOD_PATCH   = "PATCH"

	REQUEST_BODY_RAW_SYMBOL  = "@raw"
	REQUEST_BODY_FORM_SYMBOL = "@form"
)
