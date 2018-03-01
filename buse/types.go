package buse

import (
	"time"

	"github.com/itchio/butler/configurator"
	"github.com/itchio/butler/installer/bfs"
	itchio "github.com/itchio/go-itchio"
)

// must be kept in sync with clients, see for example
// https://github.com/itchio/node-butler

//----------------------------------------------------------------------
// Version
//----------------------------------------------------------------------

// Retrieves the version of the butler instance the client
// is connected to.
//
// This endpoint is meant to gather information when reporting
// issues, rather than feature sniffing. Conforming clients should
// automatically download new versions of butler, see the **Updating** section.
//
// @name Version.Get
// @category Utilities
// @tags Offline
// @caller client
type VersionGetParams struct{}

type VersionGetResult struct {
	// Something short, like `v8.0.0`
	Version string `json:"version"`

	// Something long, like `v8.0.0, built on Aug 27 2017 @ 01:13:55, ref d833cc0aeea81c236c81dffb27bc18b2b8d8b290`
	VersionString string `json:"versionString"`
}

//----------------------------------------------------------------------
// Session
//----------------------------------------------------------------------

// Lists remembered sessions
//
// @name Session.List
// @category Session
// @caller client
type SessionListParams struct {
}

type SessionListResult struct {
	// A list of remembered sessions
	Sessions []*Session `json:"sessions"`
}

// Represents a user for which we have session information,
// ie. that we can connect as, etc.
type Session struct {
	// itch.io user ID, doubling as session ID
	ID int64 `json:"id"`

	// Timestamp the user last connected at (to the client)
	LastConnected time.Time `json:"lastConnected"`

	// User information
	User *itchio.User `json:"user"`
}

// Add a new session by password login
//
// @name Session.LoginWithPassword
// @category Session
// @caller client
type SessionLoginWithPasswordParams struct {
	// The username (or e-mail) to use for login
	Username string `json:"username"`

	// The password to use
	Password string `json:"password"`
}

type SessionLoginWithPasswordResult struct {
	// Information for the new session, now remembered
	Session *Session `json:"session"`

	// Session cookie for website
	Cookie map[string]string `json:"cookie"`
}

// Ask the user to solve a captcha challenge
// Sent during @@SessionLoginWithPasswordParams if certain
// conditions are met.
//
// @name Session.RequestCaptcha
// @category Session
// @caller server
type SessionRequestCaptchaParams struct {
	// Address of page containing a recaptcha widget
	RecaptchaURL string `json:"recaptchaUrl"`
}

type SessionRequestCaptchaResult struct {
	// The response given by recaptcha after it's been filled
	RecaptchaResponse string `json:"recaptchaResponse"`
}

// Ask the user to provide a TOTP token.
// Sent during @@SessionLoginWithPasswordParams if the user has
// two-factor authentication enabled.
//
// @name Session.RequestTOTP
// @category Session
// @caller server
type SessionRequestTOTPParams struct {
}

type SessionRequestTOTPResult struct {
	// The TOTP code entered by the user
	Code string `json:"code"`
}

// Use saved login credentials to validate a session.
//
// @name Session.UseSavedLogin
// @category Session
// @caller client
type SessionUseSavedLoginParams struct {
	SessionID int64 `json:"sessionId"`
}

type SessionUseSavedLoginResult struct {
	// Information for the now validated session
	Session *Session `json:"session"`
}

// Forgets a remembered session - it won't appear in the
// @@SessionListParams results anymore.
//
// @name Session.Forget
// @category Session
// @caller client
type SessionForgetParams struct {
	SessionID int64 `json:"sessionId"`
}

type SessionForgetResult struct {
	// True if the session did exist (and was successfully forgotten)
	Success bool `json:"success"`
}

//----------------------------------------------------------------------
// Fetch
//----------------------------------------------------------------------

// Fetches information for an itch.io game.
//
// Sends @@FetchGameYieldNotification twice at most: first from cache,
// second from API if we're online.
//
// @name Fetch.Game
// @category Fetch
// @caller client
type FetchGameParams struct {
	// Session to use to fetch game
	SessionID int64 `json:"sessionId"`
	// Identifier of game to look for
	GameID int64 `json:"gameId"`
}

// Sent during @@FetchGameParams whenever a result is
// available.
//
// @name Fetch.Game.Yield
// @category Fetch
type FetchGameYieldNotification struct {
	// Current result for game fetching (from local DB, or API, etc.)
	Game *itchio.Game `json:"game"`
}

type FetchGameResult struct {
}

// Fetches information about a collection and the games it
// contains.
//
// Sends @@FetchCollectionYieldNotification.
//
// @name Fetch.Collection
// @category Fetch
// @caller client
type FetchCollectionParams struct {
	// Session to use to fetch game
	SessionID int64 `json:"sessionId"`
	// Identifier of the collection to look for
	CollectionID int64 `json:"collectionId"`
}

// Contains general info about a collection
//
// @name Fetch.Collection.Yield
// @category Fetch
type FetchCollectionYieldNotification struct {
	Collection *itchio.Collection `json:"collection"`
}

// Contains a range of games associated with a collection.
// The same ranges may be sent several times: once from
// the cache, another time from the API.
//
// @name Fetch.Collection.YieldGames
// @category Fetch
type FetchCollectionYieldGamesNotification struct {
	Offset int64             `json:"offset"`
	Total  int64             `json:"total"`
	Items  []*CollectionGame `json:"items"`
}

// Association between a @@Game and a @@Collection
// @category Fetch
type CollectionGame struct {
	Order int64        `json:"order"`
	Game  *itchio.Game `json:"game"`
}

type FetchCollectionResult struct {
}

// @name Fetch.MyCollections
// @category Fetch
// @caller client
type FetchMyCollectionsParams struct {
	// Session to use to fetch game
	SessionID int64 `json:"sessionId"`
}

// Sent during @@FetchMyCollectionsParams whenever new info is
// available.
//
// @name Fetch.MyCollections.Yield
// @category Fetch
type FetchMyCollectionsYieldNotification struct {
	Offset int64                `json:"offset"`
	Total  int64                `json:"total"`
	Items  []*CollectionSummary `json:"items"`
}

// Information about a collection + a few games
// @category Fetch
type CollectionSummary struct {
	Collection *itchio.Collection `json:"collection"`
	Items      []*CollectionGame  `json:"items"`
}

type FetchMyCollectionsResult struct {
}

// @name Fetch.MyGames
// @category Fetch
// @caller client
type FetchMyGamesParams struct {
	// Session to use to fetch game
	SessionID int64 `json:"sessionId"`
}

// @name Fetch.MyGames.Yield
// @category Fetch
type FetchMyGamesYieldNotification struct {
	Offset int64          `json:"offset"`
	Total  int64          `json:"total"`
	Items  []*itchio.Game `json:"items"`
}

type FetchMyGamesResult struct {
}

// @name Fetch.MyOwnedKeys
// @category Fetch
// @caller client
type FetchMyOwnedKeysParams struct {
	// Session to use to fetch game
	SessionID int64 `json:"sessionId"`
}

// @name Fetch.MyOwnedKeys.Yield
// @category Fetch
type FetchMyOwnedKeysYieldNotification struct {
	Offset int64                 `json:"offset"`
	Total  int64                 `json:"total"`
	Items  []*itchio.DownloadKey `json:"items"`
}

type FetchMyOwnedKeysResult struct {
}

//----------------------------------------------------------------------
// Game
//----------------------------------------------------------------------

// Finds uploads compatible with the current runtime, for a given game.
//
// @name Game.FindUploads
// @category Install
// @caller client
type GameFindUploadsParams struct {
	// Which game to find uploads for
	Game *itchio.Game `json:"game"`
	// The credentials to use to list uploads
	Credentials *GameCredentials `json:"credentials"`
}

type GameFindUploadsResult struct {
	// A list of uploads that were found to be compatible.
	Uploads []*itchio.Upload `json:"uploads"`
}

//----------------------------------------------------------------------
// Operation
//----------------------------------------------------------------------

// Start a new operation (installing or uninstalling).
//
// Can be cancelled by passing the same `ID` to @@OperationCancelParams.
//
// @name Operation.Start
// @category Install
// @tags Cancellable
// @caller client
type OperationStartParams struct {
	// A UUID, generated by the client, used for referring to the
	// task when cancelling it, for instance.
	ID string `json:"id"`

	// A folder that butler can use to store temporary files, like
	// partial downloads, checkpoint files, etc.
	StagingFolder string `json:"stagingFolder"`

	// Which operation to perform
	Operation Operation `json:"operation"`

	// Must be set if Operation is `install`
	// @optional
	InstallParams *InstallParams `json:"installParams,omitempty"`

	// Must be set if Operation is `uninstall`
	// @optional
	UninstallParams *UninstallParams `json:"uninstallParams,omitempty"`
}

type OperationStartResult struct{}

// @category Install
type Operation string

const (
	// Install a game (includes upgrades, heals, etc.)
	OperationInstall Operation = "install"
	// Uninstall a game
	OperationUninstall Operation = "uninstall"
)

// Attempt to gracefully cancel an ongoing operation.
//
// @name Operation.Cancel
// @category Install
// @caller client
type OperationCancelParams struct {
	// The UUID of the task to cancel, as passed to @@OperationStartParams
	ID string `json:"id"`
}

type OperationCancelResult struct{}

// InstallParams contains all the parameters needed to perform
// an installation for a game via @@OperationStartParams.
//
// @kind type
// @category Install
type InstallParams struct {
	// Which game to install
	Game *itchio.Game `json:"game"`

	// An absolute path where to install the game
	InstallFolder string `json:"installFolder"`

	// Which upload to install
	// @optional
	Upload *itchio.Upload `json:"upload"`

	// Which build to install
	// @optional
	Build *itchio.Build `json:"build"`

	// Which credentials to use to install the game
	Credentials *GameCredentials `json:"credentials"`

	// If true, do not run windows installers, just extract
	// whatever to the install folder.
	// @optional
	IgnoreInstallers bool `json:"ignoreInstallers,omitempty"`
}

// UninstallParams contains all the parameters needed to perform
// an uninstallation for a game via @@OperationStartParams.
//
// @kind type
// @category Install
type UninstallParams struct {
	// Absolute path of the folder butler should uninstall
	InstallFolder string `json:"installFolder"`
}

// GameCredentials contains all the credentials required to make API requests
// including the download key if any.
type GameCredentials struct {
	// Defaults to `https://itch.io`
	// @optional
	Server string `json:"server"`
	// A valid itch.io API key
	APIKey string `json:"apiKey"`
	// A download key identifier, or 0 if no download key is available
	// @optional
	DownloadKey int64 `json:"downloadKey"`
}

// Asks the user to pick between multiple available uploads
//
// @category Install
// @tags Dialog
// @caller server
type PickUploadParams struct {
	// An array of upload objects to choose from
	Uploads []*itchio.Upload `json:"uploads"`
}

type PickUploadResult struct {
	// The index (in the original array) of the upload that was picked,
	// or a negative value to cancel.
	Index int64 `json:"index"`
}

// Retrieves existing receipt information for an install
//
// @category Install
// @tags Deprecated
// @caller server
type GetReceiptParams struct {
	// muffin
}

type GetReceiptResult struct {
	Receipt *bfs.Receipt `json:"receipt"`
}

// Sent periodically during @@OperationStartParams to inform on the current state an operation.
//
// @name Operation.Progress
// @category Install
type OperationProgressNotification struct {
	// An overall progress value between 0 and 1
	Progress float64 `json:"progress"`
	// Estimated completion time for the operation, in seconds (floating)
	ETA float64 `json:"eta"`
	// Network bandwidth used, in bytes per second (floating)
	BPS float64 `json:"bps"`
}

// @category Install
type TaskReason string

const (
	// Task was started for an install operation
	TaskReasonInstall TaskReason = "install"
	// Task was started for an uninstall operation
	TaskReasonUninstall TaskReason = "uninstall"
)

// @category Install
type TaskType string

const (
	// We're fetching files from a remote server
	TaskTypeDownload TaskType = "download"
	// We're running an installer
	TaskTypeInstall TaskType = "install"
	// We're running an uninstaller
	TaskTypeUninstall TaskType = "uninstall"
	// We're applying some patches
	TaskTypeUpdate TaskType = "update"
	// We're healing from a signature and heal source
	TaskTypeHeal TaskType = "heal"
)

// Each operation is made up of one or more tasks. This notification
// is sent during @@OperationStartParams whenever a specific task starts.
//
// @category Install
type TaskStartedNotification struct {
	// Why this task was started
	Reason TaskReason `json:"reason"`
	// Is this task a download? An install?
	Type TaskType `json:"type"`
	// The game this task is dealing with
	Game *itchio.Game `json:"game"`
	// The upload this task is dealing with
	Upload *itchio.Upload `json:"upload"`
	// The build this task is dealing with (if any)
	Build *itchio.Build `json:"build,omitempty"`
	// Total size in bytes
	TotalSize int64 `json:"totalSize,omitempty"`
}

// Sent during @@OperationStartParams whenever a task succeeds for an operation.
//
// @category Install
type TaskSucceededNotification struct {
	Type TaskType `json:"type"`
	// If the task installed something, then this contains
	// info about the game, upload, build that were installed
	InstallResult *InstallResult `json:"installResult,omitempty"`
}

// What was installed by a subtask of @@OperationStartParams.
//
// See @@TaskSucceededNotification.
//
// @category Install
// @kind type
type InstallResult struct {
	// The game we installed
	Game *itchio.Game `json:"game"`
	// The upload we installed
	Upload *itchio.Upload `json:"upload"`
	// The build we installed
	// @optional
	Build *itchio.Build `json:"build"`
	// TODO: verdict ?
}

//----------------------------------------------------------------------
// CheckUpdate
//----------------------------------------------------------------------

// Looks for one or more game updates.
//
// Updates found are regularly sent via @@GameUpdateAvailableNotification, and
// then all at once in the result.
//
// @category Update
// @caller client
type CheckUpdateParams struct {
	// A list of items, each of it will be checked for updates
	Items []*CheckUpdateItem `json:"items"`
}

// @category Update
type CheckUpdateItem struct {
	// An UUID generated by the client, which allows it to map back the
	// results to its own items.
	ItemID string `json:"itemId"`
	// Timestamp of the last successful install operation
	InstalledAt string `json:"installedAt"`
	// Game for which to look for an update
	Game *itchio.Game `json:"game"`
	// Currently installed upload
	Upload *itchio.Upload `json:"upload"`
	// Currently installed build
	Build *itchio.Build `json:"build,omitempty"`
	// Credentials to use to list uploads
	Credentials *GameCredentials `json:"credentials"`
}

type CheckUpdateResult struct {
	// Any updates found (might be empty)
	Updates []*GameUpdate `json:"updates"`
	// Warnings messages logged while looking for updates
	Warnings []string `json:"warnings"`
}

// Sent during @@CheckUpdateParams, every time butler
// finds an update for a game. Can be safely ignored if displaying
// updates as they are found is not a requirement for the client.
//
// @category Update
// @tags Optional
type GameUpdateAvailableNotification struct {
	Update *GameUpdate `json:"update"`
}

// Describes an available update for a particular game install.
//
// @category Update
type GameUpdate struct {
	// Identifier originally passed in CheckUpdateItem
	ItemID string `json:"itemId"`
	// Game we found an update for
	Game *itchio.Game `json:"game"`
	// Upload to be installed
	Upload *itchio.Upload `json:"upload"`
	// Build to be installed (may be nil)
	Build *itchio.Build `json:"build"`
}

//----------------------------------------------------------------------
// Launch
//----------------------------------------------------------------------

// Attempt to launch an installed game.
//
// @category Launch
// @caller client
type LaunchParams struct {
	// The folder the game was installed to
	InstallFolder string `json:"installFolder"`
	// The itch.io game that was installed
	Game *itchio.Game `json:"game"`
	// The itch.io upload that was installed
	Upload *itchio.Upload `json:"upload"`
	// The itch.io build that was installed
	Build *itchio.Build `json:"build"`
	// The stored verdict from when the folder was last configured (can be null)
	Verdict *configurator.Verdict `json:"verdict"`

	// The directory to use to store installer files for prerequisites
	PrereqsDir string `json:"prereqsDir"`
	// Force installing all prerequisites, even if they're already marked as installed
	// @optional
	ForcePrereqs bool `json:"forcePrereqs,omitempty"`

	// Enable sandbox (regardless of manifest opt-in)
	Sandbox bool `json:"sandbox,omitempty"`

	// itch.io credentials to use for any necessary API
	// requests (prereqs downloads, subkeying, etc.)
	Credentials *GameCredentials `json:"credentials"`
}

type LaunchResult struct {
}

// Sent during @@LaunchParams, when the game is configured, prerequisites are installed
// sandbox is set up (if enabled), and the game is actually running.
//
// @category Launch
type LaunchRunningNotification struct{}

// Sent during @@LaunchParams, when the game has actually exited.
//
// @category Launch
type LaunchExitedNotification struct{}

// Sent during @@LaunchParams, ask the user to pick a manifest action to launch.
//
// See [itch app manifests](https://itch.io/docs/itch/integrating/manifest.html).
//
// @tags Dialogs
// @category Launch
// @caller server
type PickManifestActionParams struct {
	// A list of actions to pick from. Must be shown to the user in the order they're passed.
	Actions []*Action `json:"actions"`
}

type PickManifestActionResult struct {
	// Name of the action picked by user, or empty is we're aborting.
	Name string `json:"name"`
}

// Ask the client to perform a shell launch, ie. open an item
// with the operating system's default handler (File explorer).
//
// Sent during @@LaunchParams.
//
// @category Launch
// @caller server
type ShellLaunchParams struct {
	// Absolute path of item to open, e.g. `D:\\Games\\Itch\\garden\\README.txt`
	ItemPath string `json:"itemPath"`
}

type ShellLaunchResult struct {
}

// Ask the client to perform an HTML launch, ie. open an HTML5
// game, ideally in an embedded browser.
//
// Sent during @@LaunchParams.
//
// @category Launch
// @caller server
type HTMLLaunchParams struct {
	// Absolute path on disk to serve
	RootFolder string `json:"rootFolder"`
	// Path of index file, relative to root folder
	IndexPath string `json:"indexPath"`

	// Command-line arguments, to pass as `global.Itch.args`
	Args []string `json:"args"`
	// Environment variables, to pass as `global.Itch.env`
	Env map[string]string `json:"env"`
}

type HTMLLaunchResult struct {
}

// Ask the client to perform an URL launch, ie. open an address
// with the system browser or appropriate.
//
// Sent during @@LaunchParams.
//
// @category Launch
// @caller server
type URLLaunchParams struct {
	// URL to open, e.g. `https://itch.io/community`
	URL string `json:"url"`
}

type URLLaunchResult struct{}

// Ask the client to save verdict information after a reconfiguration.
//
// Sent during @@LaunchParams.
//
// @category Launch
// @tags Deprecated
// @caller server
type SaveVerdictParams struct {
	Verdict *configurator.Verdict `json:"verdict"`
}
type SaveVerdictResult struct{}

// Ask the user to allow sandbox setup. Will be followed by
// a UAC prompt (on Windows) or a pkexec dialog (on Linux) if
// the user allows.
//
// Sent during @@LaunchParams.
//
// @category Launch
// @tags Dialogs
// @caller server
type AllowSandboxSetupParams struct{}

type AllowSandboxSetupResult struct {
	// Set to true if user allowed the sandbox setup, false otherwise
	Allow bool `json:"allow"`
}

// Sent during @@LaunchParams, when some prerequisites are about to be installed.
//
// This is a good time to start showing a UI element with the state of prereq
// tasks.
//
// Updates are regularly provided via @@PrereqsTaskStateNotification.
//
// @category Launch
type PrereqsStartedNotification struct {
	// A list of prereqs that need to be tended to
	Tasks map[string]*PrereqTask `json:"tasks"`
}

// Information about a prerequisite task.
//
// @category Launch
type PrereqTask struct {
	// Full name of the prerequisite, for example: `Microsoft .NET Framework 4.6.2`
	FullName string `json:"fullName"`
	// Order of task in the list. Respect this order in the UI if you want consistent progress indicators.
	Order int `json:"order"`
}

// Current status of a prerequisite task
//
// Sent during @@LaunchParams, after @@PrereqsStartedNotification, repeatedly
// until all prereq tasks are done.
//
// @category Launch
type PrereqsTaskStateNotification struct {
	// Short name of the prerequisite task (e.g. `xna-4.0`)
	Name string `json:"name"`
	// Current status of the prereq
	Status PrereqStatus `json:"status"`
	// Value between 0 and 1 (floating)
	Progress float64 `json:"progress"`
	// ETA in seconds (floating)
	ETA float64 `json:"eta"`
	// Network bandwidth used in bytes per second (floating)
	BPS float64 `json:"bps"`
}

// @category Launch
type PrereqStatus string

const (
	// Prerequisite has not started downloading yet
	PrereqStatusPending PrereqStatus = "pending"
	// Prerequisite is currently being downloaded
	PrereqStatusDownloading PrereqStatus = "downloading"
	// Prerequisite has been downloaded and is pending installation
	PrereqStatusReady PrereqStatus = "ready"
	// Prerequisite is currently installing
	PrereqStatusInstalling PrereqStatus = "installing"
	// Prerequisite was installed (successfully or not)
	PrereqStatusDone PrereqStatus = "done"
)

// Sent during @@LaunchParams, when all prereqs have finished installing (successfully or not)
//
// After this is received, it's safe to close any UI element showing prereq task state.
//
// @category Launch
type PrereqsEndedNotification struct {
}

// Sent during @@LaunchParams, when one or more prerequisites have failed to install.
// The user may choose to proceed with the launch anyway.
//
// @category Launch
// @caller server
type PrereqsFailedParams struct {
	// Short error
	Error string `json:"error"`
	// Longer error (to include in logs)
	ErrorStack string `json:"errorStack"`
}

type PrereqsFailedResult struct {
	// Set to true if the user wants to proceed with the launch in spite of the prerequisites failure
	Continue bool `json:"continue"`
}

//----------------------------------------------------------------------
// CleanDownloads
//----------------------------------------------------------------------

// Look for folders we can clean up in various download folders.
// This finds anything that doesn't correspond to any current downloads
// we know about.
//
// @name CleanDownloads.Search
// @category Clean Downloads
// @caller client
type CleanDownloadsSearchParams struct {
	// A list of folders to scan for potential subfolders to clean up
	Roots []string `json:"roots"`
	// A list of subfolders to not consider when cleaning
	// (staging folders for in-progress downloads)
	Whitelist []string `json:"whitelist"`
}

// @category Clean Downloads
type CleanDownloadsSearchResult struct {
	// Entries we found that could use some cleaning (with path and size information)
	Entries []*CleanDownloadsEntry `json:"entries"`
}

// @category Clean Downloads
type CleanDownloadsEntry struct {
	// The complete path of the file or folder we intend to remove
	Path string `json:"path"`
	// The size of the folder or file, in bytes
	Size int64 `json:"size"`
}

// Remove the specified entries from disk, freeing up disk space.
//
// @name CleanDownloads.Apply
// @category Clean Downloads
// @caller client
type CleanDownloadsApplyParams struct {
	Entries []*CleanDownloadsEntry `json:"entries"`
}

// @category Clean Downloads
type CleanDownloadsApplyResult struct{}

//----------------------------------------------------------------------
// Misc.
//----------------------------------------------------------------------

// Sent any time butler needs to send a log message. The client should
// relay them in their own stdout / stderr, and collect them so they
// can be part of an issue report if something goes wrong.
type LogNotification struct {
	// Level of the message (`info`, `warn`, etc.)
	Level LogLevel `json:"level"`
	// Contents of the message.
	//
	// Note: logs may contain non-ASCII characters, or even emojis.
	Message string `json:"message"`
}

type LogLevel string

const (
	// Hidden from logs by default, noisy
	LogLevelDebug LogLevel = "debug"
	// Just thinking out loud
	LogLevelInfo LogLevel = "info"
	// We're continuing, but we're not thrilled about it
	LogLevelWarning LogLevel = "warning"
	// We're eventually going to fail loudly
	LogLevelError LogLevel = "error"
)

// Test request: asks butler to double a number twice.
// First by calling @@TestDoubleParams, then by
// returning the result of that call doubled.
//
// Use that to try out your JSON-RPC 2.0 over TCP implementation.
//
// @name Test.DoubleTwice
// @category Test
// @caller client
type TestDoubleTwiceParams struct {
	// The number to quadruple
	Number int64 `json:"number"`
}

// @category Test
type TestDoubleTwiceResult struct {
	// The input, quadrupled
	Number int64 `json:"number"`
}

// Test request: return a number, doubled. Implement that to
// use @@TestDoubleTwiceParams in your testing.
//
// @name Test.Double
// @category Test
// @caller server
type TestDoubleParams struct {
	// The number to double
	Number int64 `json:"number"`
}

// Result for Test.Double
type TestDoubleResult struct {
	// The number, doubled
	Number int64 `json:"number"`
}

type ItchPlatform string

const (
	ItchPlatformOSX     ItchPlatform = "osx"
	ItchPlatformWindows ItchPlatform = "windows"
	ItchPlatformLinux   ItchPlatform = "linux"
	ItchPlatformUnknown ItchPlatform = "unknown"
)

// Buse JSON-RPC 2.0 error codes
type Code int64

const (
	// An operation was cancelled gracefully
	CodeOperationCancelled Code = 499
	// An operation was aborted by the user
	CodeOperationAborted Code = 410
)

//==================================
// Manifests
//==================================

// A Manifest describes prerequisites (dependencies) and actions that
// can be taken while launching a game.
type Manifest struct {
	// Actions are a list of options to give the user when launching a game.
	Actions []*Action `json:"actions"`

	// Prereqs describe libraries or frameworks that must be installed
	// prior to launching a game
	Prereqs []*Prereq `json:"prereqs,omitempty"`
}

// An Action is a choice for the user to pick when launching a game.
//
// see https://itch.io/docs/itch/integrating/manifest.html
type Action struct {
	// human-readable or standard name
	Name string `json:"name"`

	// file path (relative to manifest or absolute), URL, etc.
	Path string `json:"path"`

	// icon name (see static/fonts/icomoon/demo.html, don't include `icon-` prefix)
	Icon string `json:"icon,omitempty"`

	// command-line arguments
	Args []string `json:"args,omitempty"`

	// sandbox opt-in
	Sandbox bool `json:"sandbox,omitempty"`

	// requested API scope
	Scope string `json:"scope,omitempty"`

	// don't redirect stdout/stderr, open in new console window
	Console bool `json:"console,omitempty"`

	// platform to restrict this action too
	Platform ItchPlatform `json:"platform,omitempty"`

	// localized action name
	Locales map[string]*ActionLocale `json:"locales,omitempty"`
}

type Prereq struct {
	// A prerequisite to be installed, see <https://itch.io/docs/itch/integrating/prereqs/> for the full list.
	Name string `json:"name"`
}

type ActionLocale struct {
	// A localized action name
	Name string `json:"name"`
}

// Dates

func FromDateTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func ToDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
