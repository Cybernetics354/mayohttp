package event

type Err error

type CheckEnvFile struct{}

type OpenTelescope struct {
	TeleType string
}

type SaveSession struct {
	Path string
}

type ReplaceCurrentSession struct {
	Path string
}

type Setup struct{}

type LoadSessionList struct{}

type LoadSession struct {
	Path string
}

type SetFieldValue struct {
	State string
	Value string
}

type RefreshState struct{}

type HideSpinner struct{}

type ShowSpinner struct{}

type RunRequest struct{}

type NextSection struct{}

type PrevSection struct{}

type RunPipe struct{}

type SetActivity string

type OpenEditor struct {
	State string
}

type OpenEnv struct{}

type OpenRequestBody struct{}

type OpenRequestHeader struct{}

type AddStack struct {
	State string
}

type PopStack struct{}

type PopStackRoot struct{}

type SetState struct {
	State string
}

type SaveInputSubmit struct{}

type SelectSessionItem struct{}

type DeleteSessionItem struct{}

type SelectCommandPallete struct{}

type SelectMethodPallete struct{}

type RunCommand struct {
	CommandId string
}

type RequestResult struct {
	Err error
	Res string
}

type PipeResult struct {
	Err error
	Res string
}

type CopyToClipboard struct{}

type RefreshSelectEnv struct{}

type SelectEnv struct{}

type RecalculateComponentSizes struct{}

