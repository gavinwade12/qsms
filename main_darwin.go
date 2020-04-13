// +build darwin

package main

func init() {
	gateways["messages"] = MessagesGateway{}
}
