Command-line arguments

Canvus client supports a number of command-line arguments that can be provided starting the executable.

The Canvus client application itself (mt-canvus-app.exe on Windows and mt-canvus-app on Ubuntu) can be found in the following locations:

Windows: %PROGRAMFILES%\MT Canvus\bin

Ubuntu: /opt/mt-canvus-3.1.0/bin

Usage:

mt-canvus-app [options] [canvas URL]

Available options:

--help
Show this help.

--version
Show version information.

--mt-canvus-config <filename>
Set the configuration file for MT Canvus.

--disable-dialogs
Disable native popup dialogs (error dialogs, licence wizards etc).

--activate <activation key>
Activate a Canvus license over the internet using an activation key.

--create-license-request <activation key>
Create an offline license activation request in the user's home directory using an activation key.

--activation-wizard
Open a license activation wizard, even if there are valid licenses available.

--audio-config <path>
Specify positional audio XML configuration file (audio-config.xml).

--config <path>
Specify touch configuration file (config.txt).

--screen <path>
Specify screen configuration file (screen.xml).

--no-crash-reporter
Disable the built-in crash reporting system and automatic crash dump upload.

--css <path>
Additional CSS file to be loaded by the application.

--scale <number>
Scale the application content (graphics coordinates) by a fixed factor. 1 means no scaling. On Windows, if the application has only one window, by default this is set to match the desktop scaling.

--cef-shared-texture-disabled
Disable web browser shared texture rendering pipeline on Windows. Might increase the browser stability on certain systems, but will slow down the browser.

Canvas url can be given to open a specific canvas on the first workspace on startup instead of the welcome screen. You can copy the URL from the Dashboard or write it manually. Examples:

canvus://example.com/0cb6639f-042b-42a3-9e97-74818d2fb09d (ssl) canvus+tcp://example.com:1234/0cb6639f-042b-42a3-9e97-74818d2fb09d (tcp) canvus://internal/0cb6639f-042b-42a3-9e97-74818d2fb09d (local canvas) https://example.com/open/0cb6639f-042b-42a3-9e97-74818d2fb09d (from Dashboard)

Commands for creating and restoring backups:

--backup
Make a backup of local canvases and exit. Uses exit code 0 on success.

--restore
Restore a local canvases backup from the path specified with --backup-path and exit. Any existing local canvases and assets will be deleted.

--backup-path <path>
Specify the folder where the backup is created to (with --backup) or restored from (with --restore). When creating a backup and --backup-path is specified, the backup is written directly to the given folder. If not specified when creating a backup, a new subfolder is created under the 'backup/root' folder specified in the configuration file and the backup is written to the subfolder. When restoring a backup, this parameter is mandatory.

--backup-skip <list>
Can be used to exclude parts of the data from the backup process. Only used with --backup command. <list> can be a comma-separated list of following tokens:

db

assets

canvus-folder

--backup-delete
Delete all files and folders in the backup-path prior to writing the new backup in there. Only used with --backup command.

The backed up Canvus folder is always the built-in Canvus folder location inside the canvus-data folder, even if you have specified another folder in 'system/canvus-folder' or don't have multi-user mode enabled and the Canvus folder in the application is your home folder. Backing up a custom Canvus folder or the home folder needs to be done separately.

Restoring a backup containing a Canvus folder will also always copy the files on top of the built-in Canvus folder location without deleting the existing files first.