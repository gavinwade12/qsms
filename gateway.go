package main

// Gateway is the interface for sending a message to a recipient
type Gateway interface {
	Send(recipient, text string) error
	RequiresConfiguration() bool
	PromptForConfiguration() error
}
