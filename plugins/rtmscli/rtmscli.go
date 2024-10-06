package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/chrlesur/orchestrator/internal/plugin"
)

type RtmsCliPlugin struct{}

func (p *RtmsCliPlugin) Name() string {
	return "rtmscli"
}

func (p *RtmsCliPlugin) Version() string {
	return "1.0.0"
}

func (p *RtmsCliPlugin) Execute(args map[string]interface{}) (interface{}, error) {
	// Construire la commande rtmscli à partir des arguments
	command := []string{"rtmscli"}
	for key, value := range args {
		command = append(command, fmt.Sprintf("--%s=%v", key, value))
	}

	// Exécuter la commande
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("rtmscli execution failed: %v, output: %s", err, string(output))
	}

	// Retourner le résultat
	return strings.TrimSpace(string(output)), nil
}

// Cette fonction est appelée par le système de plugins pour créer une instance du plugin
func New() plugin.Plugin {
	return &RtmsCliPlugin{}
}
