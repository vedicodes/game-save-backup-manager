# Game Save Backup Manager

A simple and efficient command-line tool for managing your game save backups.

![Screenshot](https://i.postimg.cc/6Q358zsw/Screenshot-2025-07-10-221256.jpg) 

## Description

This tool provides an easy way to create, restore, list, and delete backups of your game saves. It's a command-line application built with Go, designed to be fast, reliable, and easy to use.

## Features

- **Create Backups:** Easily create a backup of your game save file.
- **Restore Backups:** Restore a previously created backup.
- **List Backups:** View a list of all your available backups.
- **Delete Backups:** Remove unwanted backups.
- **Auto-Backup:** Automatically creates a backup of the current save before restoring another.
- **Configuration:** Customize the save file path and backup directory.
- **Cross-Platform:** Works on Windows, macOS, and Linux.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or higher)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/vedicodes/game-save-backup-manager.git
    cd game-save-backup-manager
    ```

2.  **Build the application:**
    ```sh
    go build
    ```

3.  **Run the application:**
    -   On Windows: `backup_manager.exe`
    -   On macOS/Linux: `./backup_manager`

## Usage

When you first run the application, it will create a `config.json` file in the same directory as the executable. You can edit this file to set your game's save file path and the directory where you want to store your backups.

The main menu provides the following options:

1.  **Create Backup:** Prompts for a backup name and creates a copy of your save file.
2.  **Restore Backup:** Shows a list of backups and lets you choose one to restore.
3.  **List Backups:** Displays all the backups in your backup directory.
4.  **Delete Backup:** Allows you to select and delete one or more backups.
5.  **Settings:** Configure the save file path, backup directory, and auto-backup feature.
6.  **Exit:** Closes the application.

## Configuration

The `config.json` file has the following structure:

```json
{
  "save_path": "path/to/your/game.sav",
  "backup_dir": "path/to/your/backups",
  "auto_backup": true
}
```

-   `save_path`: The full path to your game's save file.
-   `backup_dir`: The directory where you want to store your backups.
-   `auto_backup`: If `true`, the tool will automatically back up the current save file before restoring another.

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
