package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
)

// Config holds the CLI settings
type Config struct {
	SavePath   string `json:"save_path"`
	BackupDir  string `json:"backup_dir"`
	AutoBackup bool   `json:"auto_backup"`
}

// Backup represents a backup file
type Backup struct {
	Name      string
	Path      string
	CreatedAt time.Time
}

// Colors for CLI output
var (
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	white  = color.New(color.FgWhite).SprintFunc()
)

// Icons for CLI output
const (
	iconSuccess  = "✅"
	iconError    = "❗"
	iconDir      = "📁"
	iconRestore  = "🔄"
	iconDelete   = "🗑️"
	iconSettings = "⚙️"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s %s Error loading config: %v\n", iconError, red("ERROR:"), err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	for {
		displayMenu(config)
		choice, err := promptForChoice("Select an option (1-6)", []string{"1", "2", "3", "4", "5", "6"})
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Printf("%s %s Exiting...\n", iconError, yellow("INFO:"))
				return
			}
			fmt.Printf("%s %s Invalid input: %v\n", iconError, red("ERROR:"), err)
			continue
		}

		switch choice {
		case "1":
			createBackup(config)
		case "2":
			restoreBackup(config)
		case "3":
			listBackups(config)
		case "4":
			deleteBackups(config)
		case "5":
			config = settingsMenu(config)
		case "6":
			fmt.Printf("%s %s Thank you for using Game Save Backup Manager!\n", iconSuccess, green("INFO:"))
			fmt.Println("Press Enter to exit...")
			fmt.Scanln()
			return
		}
		fmt.Println()
	}
}

func loadConfig() (Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get executable path: %w", err)
	}
	configDir := filepath.Dir(exePath)
	configPath := filepath.Join(configDir, "config.json")

	// Default config
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		userProfile, _ = os.UserHomeDir()
	}

	config := Config{
		SavePath:   filepath.Join(userProfile, "Saved Games", "Game", "Steam ID", "game.sav"),
		BackupDir:  filepath.Join(userProfile, "Saved Games", "Game", "Steam ID", "backups"),
		AutoBackup: true,
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return Config{}, fmt.Errorf("failed to marshal default config: %w", err)
		}
		err = os.WriteFile(configPath, data, 0644)
		if err != nil {
			return Config{}, fmt.Errorf("failed to write default config: %w", err)
		}
	} else {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return Config{}, fmt.Errorf("failed to read config file: %w", err)
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			return Config{}, fmt.Errorf("failed to unmarshal config data: %w", err)
		}
	}

	if _, err := os.Stat(config.BackupDir); os.IsNotExist(err) {
		err = os.MkdirAll(config.BackupDir, 0755)
		if err != nil {
			return Config{}, fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	return config, nil
}

func saveConfig(config Config) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
}

func displayMenu(config Config) {
	clearScreen()
	fmt.Println(cyan("====================================="))
	fmt.Printf("%s %s\n", iconSettings, cyan("GAME SAVE BACKUP MANAGER"))
	fmt.Println(cyan("====================================="))
	fmt.Println()
	fmt.Printf("%s %s Current Save File: %s\n", iconDir, white("INFO:"), config.SavePath)
	fmt.Printf("%s %s Current Backup Directory: %s\n", iconDir, white("INFO:"), config.BackupDir)
	fmt.Printf("%s %s Auto-Backup on Restore: %v\n", iconSettings, white("INFO:"), config.AutoBackup)
	fmt.Println()
	fmt.Printf("1. %s Create Backup\n", iconSuccess)
	fmt.Printf("2. %s Restore Backup\n", iconRestore)
	fmt.Printf("3. %s List Backups\n", iconDir)
	fmt.Printf("4. %s Delete Backup\n", iconDelete)
	fmt.Printf("5. %s Settings\n", iconSettings)
	fmt.Printf("6. %s Exit\n", iconError)
	fmt.Println()
}

func clearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}

func promptForChoice(prompt string, validChoices []string) (string, error) {
	promptUI := promptui.Prompt{
		Label: white(prompt),
		Validate: func(input string) error {
			for _, choice := range validChoices {
				if input == choice {
					return nil
				}
			}
			return fmt.Errorf("please enter a number between 1 and %d", len(validChoices))
		},
	}
	return promptUI.Run()
}

func promptForInput(prompt string) (string, error) {
	promptUI := promptui.Prompt{
		Label: white(prompt),
	}
	result, err := promptUI.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func createBackup(config Config) {
	fmt.Println(cyan("====================================="))
	fmt.Printf("%s %s CREATE BACKUP\n", iconSuccess, cyan("CREATE BACKUP"))
	fmt.Println(cyan("====================================="))
	fmt.Println()

	if _, err := os.Stat(config.SavePath); os.IsNotExist(err) {
		fmt.Printf("%s %s Save file not found at: %s\n", iconError, red("ERROR:"), config.SavePath)
		fmt.Printf("%s %s Please check the path in Settings.\n", iconError, red("ERROR:"))
		waitForEnter()
		return
	}

	backupName, err := promptForInput("Enter backup name (press Enter for default)")
	if err != nil {
		if err != promptui.ErrInterrupt {
			fmt.Printf("%s %s Failed to read input: %v\n", iconError, red("ERROR:"), err)
		}
		waitForEnter()
		return
	}

	if backupName == "" {
		backupName = fmt.Sprintf("Backup_%s", time.Now().Format("2006-01-02_15-04-05"))
	}

	backupPath := filepath.Join(config.BackupDir, backupName+".sav")
	counter := 1
	baseName := backupName
	for {
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		backupName = fmt.Sprintf("%s_%d", baseName, counter)
		backupPath = filepath.Join(config.BackupDir, backupName+".sav")
		counter++
	}

	data, err := os.ReadFile(config.SavePath)
	if err != nil {
		fmt.Printf("%s %s Failed to read save file: %v\n", iconError, red("ERROR:"), err)
		waitForEnter()
		return
	}

	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		fmt.Printf("%s %s Failed to create backup: %v\n", iconError, red("ERROR:"), err)
	} else {
		createdAt, _ := getFileCreationTime(backupPath)
		fmt.Printf("%s %s Backup created successfully!\n", iconSuccess, green("SUCCESS:"))
		fmt.Printf("%s %s Backup name: %s\n", iconSuccess, green("INFO:"), backupName)
		fmt.Printf("%s %s Created at: %s\n", iconSuccess, green("INFO:"), createdAt.Format("01/02/2006 03:04:05 PM"))
	}

	waitForEnter()
}

func restoreBackup(config Config) {
	fmt.Println(cyan("====================================="))
	fmt.Printf("%s %s RESTORE BACKUP\n", iconRestore, cyan("RESTORE BACKUP"))
	fmt.Println(cyan("====================================="))
	fmt.Println()

	backups, err := listBackupsInternal(config)
	if err != nil {
		fmt.Printf("%s %s Failed to list backups: %v\n", iconError, red("ERROR:"), err)
		waitForEnter()
		return
	}

	if len(backups) == 0 {
		fmt.Printf("%s %s No backups found.\n", iconError, red("INFO:"))
		waitForEnter()
		return
	}

	items := make([]string, len(backups))
	for i, backup := range backups {
		items[i] = fmt.Sprintf("%s (Created: %s)", backup.Name, backup.CreatedAt.Format("01/02/2006 03:04:05 PM"))
	}

	prompt := promptui.Select{
		Label: white("Select a backup to restore (or cancel)"),
		Items: append(items, "Cancel"),
	}
	index, _, err := prompt.Run()
	if err != nil {
		if err != promptui.ErrInterrupt {
			fmt.Printf("%s %s Failed to select backup: %v\n", iconError, red("ERROR:"), err)
		}
		waitForEnter()
		return
	}
	if index == len(items) {
		fmt.Printf("%s %s Restore cancelled.\n", iconError, yellow("INFO:"))
		waitForEnter()
		return
	}

	selectedBackup := backups[index]
	fmt.Println()
	fmt.Printf("%s %s WARNING: This will overwrite your current save file!\n", iconError, yellow("WARNING:"))
	fmt.Printf("%s %s Selected backup: %s\n", iconRestore, yellow("INFO:"), selectedBackup.Name)
	fmt.Println()

	confirm, err := promptForInput("Are you sure you want to restore this backup? (y/N)")
	if err != nil || strings.ToLower(confirm) != "y" {
		fmt.Printf("%s %s Restore cancelled.\n", iconError, yellow("INFO:"))
		waitForEnter()
		return
	}

	if config.AutoBackup {
		if _, err := os.Stat(config.SavePath); !os.IsNotExist(err) {
			autoBackupName := fmt.Sprintf("AutoBackup_%s", time.Now().Format("2006-01-02_15-04-05"))
			autoBackupPath := filepath.Join(config.BackupDir, autoBackupName+".sav")
			data, err := os.ReadFile(config.SavePath)
			if err == nil {
				err = os.WriteFile(autoBackupPath, data, 0644)
				if err == nil {
					fmt.Printf("%s %s Auto-backup of current save created: %s\n", iconSuccess, green("SUCCESS:"), autoBackupName)
				}
			}
		}
	}

	data, err := os.ReadFile(selectedBackup.Path)
	if err != nil {
		fmt.Printf("%s %s Failed to read backup: %v\n", iconError, red("ERROR:"), err)
	} else {
		err = os.WriteFile(config.SavePath, data, 0644)
		if err != nil {
			fmt.Printf("%s %s Failed to restore backup: %v\n", iconError, red("ERROR:"), err)
		} else {
			fmt.Printf("%s %s Backup restored successfully!\n", iconSuccess, green("SUCCESS:"))
		}
	}

	waitForEnter()
}

func listBackups(config Config) {
	fmt.Println(cyan("====================================="))
	fmt.Printf("%s %s BACKUP LIST\n", iconDir, cyan("BACKUP LIST"))
	fmt.Println(cyan("====================================="))
	fmt.Println()

	backups, err := listBackupsInternal(config)
	if err != nil {
		fmt.Printf("%s %s Failed to list backups: %v\n", iconError, red("ERROR:"), err)
	} else if len(backups) == 0 {
		fmt.Printf("%s %s No backups found.\n", iconError, red("INFO:"))
	} else {
		for i, backup := range backups {
			fmt.Printf("%d. %s %s (Created: %s)\n", i+1, iconDir, white(backup.Name), backup.CreatedAt.Format("01/02/2006 03:04:05 PM"))
		}
	}

	waitForEnter()
}

func listBackupsInternal(config Config) ([]Backup, error) {
	files, err := os.ReadDir(config.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []Backup
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sav") {
			path := filepath.Join(config.BackupDir, file.Name())
			createdAt, err := getFileCreationTime(path)
			if err != nil {
				// Log error or handle it, for now, skip the file
				continue
			}
			name := strings.TrimSuffix(file.Name(), ".sav")
			backups = append(backups, Backup{
				Name:      name,
				Path:      path,
				CreatedAt: createdAt,
			})
		}
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

func deleteBackups(config Config) {
	fmt.Println(cyan("====================================="))
	fmt.Printf("%s %s DELETE BACKUP\n", iconDelete, cyan("DELETE BACKUP"))
	fmt.Println(cyan("====================================="))
	fmt.Println()

	backups, err := listBackupsInternal(config)
	if err != nil {
		fmt.Printf("%s %s Failed to list backups: %v\n", iconError, red("ERROR:"), err)
		waitForEnter()
		return
	}

	if len(backups) == 0 {
		fmt.Printf("%s %s No backups found.\n", iconError, red("INFO:"))
		waitForEnter()
		return
	}

	items := make([]string, len(backups))
	for i, backup := range backups {
		items[i] = fmt.Sprintf("%s (Created: %s)", backup.Name, backup.CreatedAt.Format("01/02/2006 03:04:05 PM"))
	}

	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message: "Select backups to delete (use space to select, enter to confirm):",
		Options: items,
	}
	err = survey.AskOne(prompt, &selectedIndices)
	if err != nil {
		fmt.Printf("%s %s Deletion cancelled.\n", iconError, yellow("INFO:"))
		waitForEnter()
		return
	}

	if len(selectedIndices) == 0 {
		fmt.Printf("%s %s No backups selected.\n", iconError, yellow("INFO:"))
		waitForEnter()
		return
	}

	fmt.Println()
	fmt.Printf("%s %s WARNING: This will permanently delete the selected backups!\n", iconError, yellow("WARNING:"))
	for _, index := range selectedIndices {
		backupName := strings.Split(items[index], " (Created:")[0]
		fmt.Printf(" - %s %s\n", iconDelete, yellow(backupName))
	}
	fmt.Println()

	confirm, err := promptForInput("Are you sure? (y/N)")
	if err != nil || strings.ToLower(confirm) != "y" {
		fmt.Printf("%s %s Deletion cancelled.\n", iconError, yellow("INFO:"))
		waitForEnter()
		return
	}

	deletedCount := 0
	for _, index := range selectedIndices {
		backup := backups[index]
		err := os.Remove(backup.Path)
		if err != nil {
			fmt.Printf("%s %s Failed to delete %s: %v\n", iconError, red("ERROR:"), backup.Name, err)
		} else {
			deletedCount++
		}
	}

	if deletedCount > 0 {
		fmt.Printf("%s %s %d backup(s) deleted successfully!\n", iconSuccess, green("SUCCESS:"), deletedCount)
	}
	waitForEnter()
}

func settingsMenu(config Config) Config {
	for {
		fmt.Println(cyan("====================================="))
		fmt.Printf("%s %s SETTINGS\n", iconSettings, cyan("SETTINGS"))
		fmt.Println(cyan("====================================="))
		fmt.Println()
		fmt.Printf("%s %s Current Save File Path: %s\n", iconDir, white("INFO:"), config.SavePath)
		fmt.Printf("%s %s Current Backup Directory: %s\n", iconDir, white("INFO:"), config.BackupDir)
		fmt.Printf("%s %s Auto-Backup on Restore: %v\n", iconSettings, white("INFO:"), config.AutoBackup)
		fmt.Println()
		fmt.Printf("1. %s Change Save File Path\n", iconSettings)
		fmt.Printf("2. %s Change Backup Directory\n", iconSettings)
		fmt.Printf("3. %s Toggle Auto-Backup on Restore\n", iconSettings)
		fmt.Printf("4. %s Test Save File Path\n", iconSettings)
		fmt.Printf("5. %s Open Backup Directory\n", iconDir)
		fmt.Printf("6. %s Back to Main Menu\n", iconSuccess)
		fmt.Println()

		choice, err := promptForChoice("Select an option (1-6)", []string{"1", "2", "3", "4", "5", "6"})
		if err != nil {
			if err == promptui.ErrInterrupt {
				return config // Exit settings on interrupt
			}
			fmt.Printf("%s %s Invalid input: %v\n", iconError, red("ERROR:"), err)
			continue
		}

		switch choice {
		case "1":
			fmt.Println()
			fmt.Printf("%s %s Current path: %s\n", iconDir, white("INFO:"), config.SavePath)
			newPath, err := promptForInput("Enter new save file path")
			if err == nil && newPath != "" {
				config.SavePath = newPath
				if err := saveConfig(config); err != nil {
					fmt.Printf("%s %s Failed to save config: %v\n", iconError, red("ERROR:"), err)
				}
			}
		case "2":
			fmt.Println()
			fmt.Printf("%s %s Current directory: %s\n", iconDir, white("INFO:"), config.BackupDir)
			newDir, err := promptForInput("Enter new backup directory")
			if err == nil && newDir != "" {
				config.BackupDir = newDir
				if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
					fmt.Printf("%s %s Failed to create backup directory: %v\n", iconError, red("ERROR:"), err)
				}
				if err := saveConfig(config); err != nil {
					fmt.Printf("%s %s Failed to save config: %v\n", iconError, red("ERROR:"), err)
				}
			}
		case "3":
			fmt.Println()
			config.AutoBackup = !config.AutoBackup
			status := "DISABLED"
			if config.AutoBackup {
				status = "ENABLED"
			}
			fmt.Printf("%s %s Auto-backup has been %s\n", iconSuccess, green("SUCCESS:"), status)
			if err := saveConfig(config); err != nil {
				fmt.Printf("%s %s Failed to save config: %v\n", iconError, red("ERROR:"), err)
			}
			waitForEnter()
		case "4":
			fmt.Println()
			if _, err := os.Stat(config.SavePath); os.IsNotExist(err) {
				fmt.Printf("%s %s Save file not found at: %s\n", iconError, red("ERROR:"), config.SavePath)
			} else {
				fmt.Printf("%s %s Save file found at: %s\n", iconSuccess, green("SUCCESS:"), config.SavePath)
			}
			waitForEnter()
		case "5":
			openExplorer(config.BackupDir)
			waitForEnter()
		case "6":
			return config
		}
		fmt.Println()
	}
}

func openExplorer(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("%s %s Failed to open directory: %v\n", iconError, red("ERROR:"), err)
	} else {
		fmt.Printf("%s %s Directory opened.\n", iconSuccess, green("SUCCESS:"))
	}
}

func waitForEnter() {
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
}
