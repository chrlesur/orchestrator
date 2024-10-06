package ui

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/chrlesur/orchestrator/internal/job"
	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/internal/pipeline"
	"github.com/chrlesur/orchestrator/internal/plugin"
	"github.com/chrlesur/orchestrator/pkg/logger"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	app              *tview.Application
	jobManager       *job.Manager
	pipelineManager  *pipeline.Manager
	pluginManager    *plugin.PluginManager
	jobList          *tview.List
	pipelineList     *tview.List
	logView          *tview.TextView
	detailView       *tview.TextView
	statsView        *tview.TextView
	inputField       *tview.InputField
	logLevel         logger.LogLevel
	logLevelDropDown *tview.DropDown
}

func NewTUI(jobManager *job.Manager, pipelineManager *pipeline.Manager, pluginManager *plugin.PluginManager) *TUI {
	t := &TUI{
		app:             tview.NewApplication(),
		jobManager:      jobManager,
		pipelineManager: pipelineManager,
		pluginManager:   pluginManager,
		jobList:         tview.NewList().ShowSecondaryText(false),
		pipelineList:    tview.NewList().ShowSecondaryText(false),
		logView:         tview.NewTextView().SetDynamicColors(true),
		detailView:      tview.NewTextView().SetDynamicColors(true),
		statsView:       tview.NewTextView().SetDynamicColors(true),
		inputField:      tview.NewInputField().SetLabel("Command: "),
		logLevel:        logger.INFO,
	}

	t.logView.SetChangedFunc(func() { t.app.Draw() })
	t.setupUI()
	return t
}

func (t *TUI) setupUI() {
	// Configuration des différentes vues
	t.jobList.SetTitle("Jobs").SetBorder(true).SetTitleColor(tcell.ColorGreen)
	t.pipelineList.SetTitle("Pipelines").SetBorder(true).SetTitleColor(tcell.ColorBlue)
	t.logView.SetTitle("Logs").SetBorder(true).SetTitleColor(tcell.ColorYellow)
	t.detailView.SetTitle("Details").SetBorder(true).SetTitleColor(tcell.GetColor("#FF00FF"))   // Magenta
	t.statsView.SetTitle("Statistics").SetBorder(true).SetTitleColor(tcell.GetColor("#00FFFF")) // Cyan

	// Configuration des fonctions de sélection
	t.jobList.SetSelectedFunc(t.showJobDetails)
	t.pipelineList.SetSelectedFunc(t.showPipelineDetails)

	// Configuration de l'input field
	t.inputField.SetLabel("Command: ").SetFieldWidth(0)
	t.inputField.SetDoneFunc(t.handleCommand)
	t.inputField.SetChangedFunc(func(text string) {
		if text == "" {
			t.showHelp()
		}
	})

	// Configuration du menu déroulant pour le niveau de log
	t.logLevelDropDown = tview.NewDropDown().
		SetLabel("Log Level: ").
		SetOptions([]string{"DEBUG", "INFO", "WARNING", "ERROR"}, t.setLogLevel)
	t.updateLogLevelDisplay() // Initialise l'affichage du niveau de log

	// Configuration de la liste des plugins
	pluginMenu := tview.NewList().ShowSecondaryText(false)
	for _, pluginName := range t.pluginManager.GetLoadedPlugins() {
		pluginMenu.AddItem(pluginName, "", 0, nil)
	}
	pluginMenu.SetTitle("Plugins").SetBorder(true)

	// Organisation de l'interface
	leftPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.jobList, 0, 4, false).
		AddItem(t.pipelineList, 0, 4, false).
		AddItem(pluginMenu, 0, 1, false)

	rightPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.detailView, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(t.logLevelDropDown, 1, 0, false).
			AddItem(t.logView, 0, 7, false), 0, 7, false).
		AddItem(t.statsView, 0, 1, false)

	mainFlex := tview.NewFlex().
		AddItem(leftPane, 0, 3, false).
		AddItem(rightPane, 0, 2, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, false).
		AddItem(t.inputField, 1, 0, true)

	// Configuration de l'application
	t.app.SetRoot(root, true).SetFocus(t.inputField)

	// Ajout d'un gestionnaire d'événements global pour les touches
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Changement de focus entre les éléments principaux
			if t.app.GetFocus() == t.inputField {
				t.app.SetFocus(t.jobList)
			} else if t.app.GetFocus() == t.jobList {
				t.app.SetFocus(t.pipelineList)
			} else {
				t.app.SetFocus(t.inputField)
			}
			return nil
		case tcell.KeyCtrlC:
			// Quitter l'application
			t.app.Stop()
			return nil
		}
		return event
	})

	// Message d'accueil
	t.detailView.SetText("Welcome to the Orchestrator. Type 'help' to see available commands.")
	logger.Info(fmt.Sprintf("Welcome to the Orchestrator"))
}

func (t *TUI) updateJobList() {
	t.jobList.Clear()
	for _, job := range t.jobManager.GetJobs() {
		t.jobList.AddItem(fmt.Sprintf("%s - %s (%s)", job.Name, job.ID, job.Status), "", 0, nil)
	}
}

func (t *TUI) updatePipelineList() {
	t.pipelineList.Clear()
	for _, pipeline := range t.pipelineManager.GetPipelines() {
		t.pipelineList.AddItem(fmt.Sprintf("%s - %s", pipeline.ID, pipeline.Status), "", 0, nil)
	}
}

func (t *TUI) showJobDetails(index int, mainText string, secondaryText string, shortcut rune) {
	parts := strings.SplitN(mainText, " - ", 2)
	if len(parts) < 2 {
		return
	}
	jobID := strings.TrimSpace(strings.Split(parts[1], " ")[0])
	job, err := t.jobManager.GetJob(jobID)
	if err != nil {
		t.detailView.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	details := fmt.Sprintf("Job Name: %s\nJob ID: %s\nStatus: %s\nCommand: %s\nArgs: %v\nStart Time: %s\nEnd Time: %s\nResult: %s\nError: %v",
		job.Name, job.ID, job.Status, job.Command, job.Args, job.StartTime, job.EndTime, job.Result, job.Error)
	t.detailView.SetText(details)
}

func (t *TUI) showPipelineDetails(index int, mainText string, secondaryText string, shortcut rune) {
	pipelineID := strings.Split(mainText, " - ")[0]
	pipeline, err := t.pipelineManager.GetPipeline(pipelineID)
	if err != nil {
		t.detailView.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	details := fmt.Sprintf("Pipeline ID: %s\nName: %s\nStatus: %s\nStart Time: %s\nEnd Time: %s\nJobs: %d\nScheduled At: %s",
		pipeline.ID, pipeline.Name, pipeline.Status, pipeline.StartTime, pipeline.EndTime, len(pipeline.Jobs), pipeline.ScheduledAt)
	t.detailView.SetText(details)
}

func (t *TUI) handleCommand(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	cmd := t.inputField.GetText()
	t.inputField.SetText("")

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "help":
		t.showHelp()
	case "addjob":
		t.handleAddJob(parts[1:])
	case "addpipeline":
		t.handleAddPipeline(parts[1:])
	case "executeplugin":
		t.handleExecutePlugin(parts[1:])
	case "setloglevel":
		t.handleSetLogLevel(parts[1:])
	default:
		logger.Info(fmt.Sprintf("Unknown command: %s. Type 'help' for available commands.", parts[0]))
		t.detailView.SetText(fmt.Sprintf("Unknown command: %s. Type 'help' for available commands.", parts[0]))
	}
}

func (t *TUI) handleAddJob(args []string) {
	if len(args) < 3 {
		t.detailView.SetText("Usage: addjob <name> <command> <arg1> <arg2> ...")
		return
	}

	name := args[0]
	command := args[1]
	jobArgs := args[2:]

	job, err := t.jobManager.CreateJob(name, command, jobArgs, "")
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding job: %v", err))
		t.detailView.SetText(fmt.Sprintf("Error adding job: %v", err))
	} else {
		logger.Info(fmt.Sprintf("Job added: %s (ID: %s)", name, job.ID))
		t.detailView.SetText(fmt.Sprintf("Job added successfully: %s (ID: %s)", name, job.ID))
		t.updateJobList()
	}
}

func (t *TUI) handleAddPipeline(args []string) {
	if len(args) < 2 {
		logger.Info("Usage: addpipeline <id> <name> <job1> <job2> ...")
		return
	}

	id := args[0]
	name := args[1]
	jobIDs := args[2:]

	jobs := make([]*models.Job, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		j, err := t.jobManager.GetJob(jobID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting job %s: %v", jobID, err))
			return
		}
		jobs = append(jobs, j)
	}

	newPipeline := &models.Pipeline{
		ID:          id,
		Name:        name,
		Jobs:        jobs,
		Status:      models.PipelineStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Minute),
	}
	err := t.pipelineManager.AddPipeline(newPipeline)
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding pipeline: %v", err))
	} else {
		logger.Info(fmt.Sprintf("Pipeline added: %s", id))
		t.updatePipelineList()
	}
}

func (t *TUI) handleExecutePlugin(args []string) {
	if len(args) < 2 {
		logger.Info("Usage: executeplugin <plugin_name> <arg1> <arg2> ...")
		return
	}

	pluginName := args[0]
	pluginArgs := make(map[string]interface{})
	for i, arg := range args[1:] {
		pluginArgs[fmt.Sprintf("arg%d", i+1)] = arg
	}

	result, err := t.pluginManager.ExecutePlugin(pluginName, pluginArgs)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to execute plugin %s: %v", pluginName, err))
	} else {
		logger.Info(fmt.Sprintf("Plugin %s executed successfully. Result: %v", pluginName, result))
	}
}

func (t *TUI) showHelp() {
	helpText := `Available commands:
    help - Display this help message
    addjob <name> <command> <arg1> <arg2> ... - Add a new job
    addpipeline <id> <name> <job1> <job2> ... - Add a new pipeline
    executeplugin <plugin_name> <arg1> <arg2> ... - Execute a plugin
    setloglevel <DEBUG|INFO|WARNING|ERROR> - Set the log level`

	t.detailView.SetText(helpText)
}

func (t *TUI) setLogLevel(text string, index int) {
	var newLevel logger.LogLevel
	switch text {
	case "DEBUG":
		newLevel = logger.DEBUG
	case "INFO":
		newLevel = logger.INFO
	case "WARNING":
		newLevel = logger.WARNING
	case "ERROR":
		newLevel = logger.ERROR
	default:
		return
	}

	logger.SetLogLevel(newLevel)
}

func (t *TUI) updateLogLevelDisplay() {
	currentLevel := logger.GetCurrentLogLevel()
	levelString := logger.GetLevelString(currentLevel)
	t.logLevelDropDown.SetLabel(fmt.Sprintf("Log Level (%s): ", levelString))
}

func (t *TUI) updateLogs() {
	logs := logger.GetLogs(t.logLevel)
	t.logView.Clear()
	for _, log := range logs {
		t.logView.Write([]byte(fmt.Sprintf("[%s] [%s] %s\n",
			log.Timestamp.Format("15:04:05"),
			strings.ToLower(logger.GetLevelString(log.Level)),
			log.Message)))
	}
}

func (t *TUI) updateStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := fmt.Sprintf(
		"CPU Usage: %.2f%%\n"+
			"Memory Usage: %v MB\n"+
			"Goroutines: %d\n"+
			"Jobs Running: %d\n"+
			"Pipelines Running: %d",
		t.getCPUUsage(),
		m.Alloc/1024/1024,
		runtime.NumGoroutine(),
		t.getRunningJobsCount(),
		t.getRunningPipelinesCount(),
	)

	t.statsView.Clear()
	t.statsView.SetText(stats)
}

func (t *TUI) getCPUUsage() float64 {
	// Cette fonction est une simplification. Pour obtenir l'utilisation réelle du CPU,
	// vous devriez utiliser une bibliothèque comme github.com/shirou/gopsutil
	return 0.0
}

func (t *TUI) getRunningJobsCount() int {
	count := 0
	for _, job := range t.jobManager.GetJobs() {
		if job.Status == models.JobStatusRunning {
			count++
		}
	}
	return count
}

func (t *TUI) getRunningPipelinesCount() int {
	count := 0
	for _, pipeline := range t.pipelineManager.GetPipelines() {
		if pipeline.Status == models.PipelineStatusRunning {
			count++
		}
	}
	return count
}

func (t *TUI) Run() error {
	// Mise à jour périodique des listes, des logs et des statistiques
	go func() {
		for {
			t.app.QueueUpdateDraw(func() {
				t.updateJobList()
				t.updatePipelineList()
				t.updateLogs()
				t.updateStats()
				t.updateLogLevelDisplay()
			})
			time.Sleep(5 * time.Second)
		}
	}()

	return t.app.Run()
}

func (t *TUI) handleSetLogLevel(args []string) {
	if len(args) != 1 {
		t.detailView.SetText("Usage: setloglevel <DEBUG|INFO|WARNING|ERROR>")
		return
	}

	level := strings.ToUpper(args[0])
	var newLevel logger.LogLevel
	switch level {
	case "DEBUG":
		newLevel = logger.DEBUG
	case "INFO":
		newLevel = logger.INFO
	case "WARNING":
		newLevel = logger.WARNING
	case "ERROR":
		newLevel = logger.ERROR
	default:
		t.detailView.SetText("Invalid log level. Use DEBUG, INFO, WARNING, or ERROR.")
		return
	}

	logger.SetLogLevel(newLevel)
	t.updateLogLevelDisplay()
	t.detailView.SetText(fmt.Sprintf("Log level set to %s", level))
}
