package itchio

// User represents an itch.io account, with basic profile info
type User struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`

	// The user's username (used for login)
	Username string `json:"username"`
	// The user's display name: human-friendly, may contain spaces, unicode etc.
	DisplayName string `json:"displayName"`

	// Has the user opted into creating games?
	Developer bool `json:"developer"`
	// Is the user part of itch.io's press program?
	PressUser bool `json:"pressUser"`

	// The address of the user's page on itch.io
	URL string `json:"url"`
	// User's avatar, may be a GIF
	CoverURL string `json:"coverUrl"`
	// Static version of user's avatar, only set if the main cover URL is a GIF
	StillCoverURL string `json:"stillCoverUrl"`
}

// Game represents a page on itch.io, it could be a game,
// a tool, a comic, etc.
type Game struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`
	// Canonical address of the game's page on itch.io
	URL string `json:"url"`

	// Human-friendly title (may contain any character)
	Title string `json:"title"`
	// Human-friendly short description
	ShortText string `json:"shortText"`
	// Downloadable game, html game, etc.
	Type GameType `json:"type"`
	// Classification: game, tool, comic, etc.
	Classification GameClassification `json:"classification"`

	// Configuration for embedded (HTML5) games
	// @optional
	Embed *GameEmbedInfo `json:"embed"`

	// Cover url (might be a GIF)
	CoverURL string `json:"coverUrl"`
	// Non-gif cover url, only set if main cover url is a GIF
	StillCoverURL string `json:"stillCoverUrl"`

	// Date the game was created
	CreatedAt string `json:"createdAt"`
	// Date the game was published, empty if not currently published
	PublishedAt string `json:"publishedAt"`

	// Price in cents of a dollar
	MinPrice int64 `json:"minPrice"`
	// Can this game be bought?
	CanBeBought bool `json:"canBeBought"`
	// Is this game downloadable by press users for free?
	InPressSystem bool `json:"inPressSystem"`
	// Does this game have a demo that can be downloaded for free?
	HasDemo bool `json:"hasDemo"`

	// Does this game have an upload tagged with 'macOS compatible'? (creator-controlled)
	OSX bool `json:"pOsx"`
	// Does this game have an upload tagged with 'Linux compatible'? (creator-controlled)
	Linux bool `json:"pLinux"`
	// Does this game have an upload tagged with 'Windows compatible'? (creator-controlled)
	Windows bool `json:"pWindows"`
	// Does this game have an upload tagged with 'Android compatible'? (creator-controlled)
	Android bool `json:"pAndroid"`

	// The user account this game is associated to
	// @optional
	User *User `json:"user" gorm:"-"`

	// ID of the user account this game is associated to
	UserID int64 `json:"userId"`

	// The best current sale for this game
	// @optional
	Sale *Sale `json:"sale" gorm:"-"`
}

// Type of an itch.io game page, mostly related to
// how it should be presented on web (downloadable or embed)
type GameType string

const (
	// downloadable
	GameTypeDefault GameType = "default"
	// .swf (legacy)
	GameTypeFlash GameType = "flash"
	// .unity3d (legacy)
	GameTypeUnity GameType = "unity"
	// .jar (legacy)
	GameTypeJava GameType = "java"
	// .html (thriving)
	GameTypeHTML GameType = "html"
)

// Creator-picked classification for a page
type GameClassification string

const (
	// something you can play
	GameClassificationGame GameClassification = "game"
	// all software pretty much
	GameClassificationTool GameClassification = "tool"
	// assets: graphics, sounds, etc.
	GameClassificationAssets GameClassification = "assets"
	// game mod (no link to game, purely creator tagging)
	GameClassificationGameMod GameClassification = "game_mod"
	// printable / board / card game
	GameClassificationPhysicalGame GameClassification = "physical_game"
	// bunch of music files
	GameClassificationSoundtrack GameClassification = "soundtrack"
	// anything that creators think don't fit in any other category
	GameClassificationOther GameClassification = "other"
	// comic book (pdf, jpg, specific comic formats, etc.)
	GameClassificationComic GameClassification = "comic"
	// book (pdf, jpg, specific e-book formats, etc.)
	GameClassificationBook GameClassification = "book"
)

// Presentation information for embed games
type GameEmbedInfo struct {
	// width of the initial viewport, in pixels
	Width int64 `json:"width"`

	// height of the initial viewport, in pixels
	Height int64 `json:"height"`

	// for itch.io website, whether or not a fullscreen button should be shown
	Fullscreen bool `json:"fullscreen"`
}

// Describes a discount for a game.
type Sale struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`
	// Discount rate in percent.
	// Can be negative, see https://itch.io/updates/introducing-reverse-sales
	Rate float64 `json:"rate"`
	// Timestamp the sale started at
	StartDate string `json:"startDate"`
	// Timestamp the sale ends at
	EndDate string `json:"endDate"`
}

// An Upload is a downloadable file. Some are wharf-enabled, which means
// they're actually a "channel" that may contain multiple builds, pushed
// with <https://github.com/itchio/butler>
type Upload struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`
	// Original file name (example: `Overland_x64.zip`)
	Filename string `json:"filename"`
	// Human-friendly name set by developer (example: `Overland for Windows 64-bit`)
	DisplayName string `json:"displayName"`
	// Size of upload in bytes. For wharf-enabled uploads, it's the archive size.
	Size int64 `json:"size"`
	// Name of the wharf channel for this upload, if it's a wharf-enabled upload
	ChannelName string `json:"channelName"`
	// Latest build for this upload, if it's a wharf-enabled upload
	Build *Build `json:"build"`
	// Is this upload a demo that can be downloaded for free?
	Demo bool `json:"demo"`
	// Is this upload a pre-order placeholder?
	Preorder bool `json:"preorder"`

	// Upload type: default, soundtrack, etc.
	Type string `json:"type"`

	// Is this upload tagged with 'macOS compatible'? (creator-controlled)
	OSX bool `json:"pOsx"`
	// Is this upload tagged with 'Linux compatible'? (creator-controlled)
	Linux bool `json:"pLinux"`
	// Is this upload tagged with 'Windows compatible'? (creator-controlled)
	Windows bool `json:"pWindows"`
	// Is this upload tagged with 'Android compatible'? (creator-controlled)
	Android bool `json:"pAndroid"`

	// Date this upload was created at
	CreatedAt string `json:"createdAt"`
	// Date this upload was last updated at (order changed, display name set, etc.)
	UpdatedAt string `json:"updatedAt"`
}

// A Collection is a set of games, curated by humans.
type Collection struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`

	// Human-friendly title for collection, for example `Couch coop games`
	Title string `json:"title"`

	// Date this collection was created at
	CreatedAt string `json:"createdAt"`
	// Date this collection was last updated at (item added, title set, etc.)
	UpdatedAt string `json:"updatedAt"`

	// Number of games in the collection. This might not be accurate
	// as some games might not be accessible to whoever is asking (project
	// page deleted, visibility level changed, etc.)
	GamesCount int64 `json:"gamesCount"`
}

// A download key is often generated when a purchase is made, it
// allows downloading uploads for a game that are not available
// for free.
type DownloadKey struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`

	// Identifier of the game to which this download key grants access
	GameID int64 `json:"gameId"`

	// Game to which this download key grants access
	Game *Game `json:"game,omitempty" gorm:"-"`

	// Date this key was created at (often coincides with purchase time)
	CreatedAt string `json:"createdAt"`
	// Date this key was last updated at
	UpdatedAt string `json:"updatedAt"`

	// Identifier of the itch.io user to which this key belongs
	OwnerID int64 `json:"ownerId"`
}

// Build contains information about a specific build
type Build struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`
	// Identifier of the build before this one on the same channel,
	// or 0 if this is the initial build.
	ParentBuildID int64 `json:"parentBuildId"`
	// State of the build: started, processing, etc.
	State BuildState `json:"state"`

	// Automatically-incremented version number, starting with 1
	Version int64 `json:"version"`
	// Value specified by developer with `--userversion` when pushing a build
	// Might not be unique across builds of a given channel.
	UserVersion string `json:"userVersion"`

	// Files associated with this build - often at least an archive,
	// a signature, and a patch. Some might be missing while the build
	// is still processing or if processing has failed.
	Files []*BuildFile `json:"files"`

	// User who pushed the build
	User User `json:"user"`
	// Timestamp the build was created at
	CreatedAt string `json:"createdAt"`
	// Timestamp the build was last updated at
	UpdatedAt string `json:"updatedAt"`
}

// BuildState describes the state of a build, relative to its initial upload, and
// its processing.
type BuildState string

const (
	// BuildStateStarted is the state of a build from its creation until the initial upload is complete
	BuildStateStarted BuildState = "started"
	// BuildStateProcessing is the state of a build from the initial upload's completion to its fully-processed state.
	// This state does not mean the build is actually being processed right now, it's just queued for processing.
	BuildStateProcessing BuildState = "processing"
	// BuildStateCompleted means the build was successfully processed. Its patch hasn't necessarily been
	// rediff'd yet, but we have the holy (patch,signature,archive) trinity.
	BuildStateCompleted BuildState = "completed"
	// BuildStateFailed means something went wrong with the build. A failing build will not update the channel
	// head and can be requeued by the itch.io team, although if a new build is pushed before they do,
	// that new build will "win".
	BuildStateFailed BuildState = "failed"
)

// BuildFile contains information about a build's "file", which could be its
// archive, its signature, its patch, etc.
type BuildFile struct {
	// Site-wide unique identifier generated by itch.io
	ID int64 `json:"id"`
	// Size of this build file
	Size int64 `json:"size"`
	// State of this file: created, uploading, uploaded, etc.
	State BuildFileState `json:"state"`
	// Type of this build file: archive, signature, patch, etc.
	Type BuildFileType `json:"type"`
	// Subtype of this build file, usually indicates compression
	SubType BuildFileSubType `json:"subType"`

	// Date this build file was created at
	CreatedAt string `json:"createdAt"`
	// Date this build file was last updated at
	UpdatedAt string `json:"updatedAt"`
}

// BuildFileState describes the state of a specific file for a build
type BuildFileState string

const (
	// BuildFileStateCreated means the file entry exists on itch.io
	BuildFileStateCreated BuildFileState = "created"
	// BuildFileStateUploading means the file is currently being uploaded to storage
	BuildFileStateUploading BuildFileState = "uploading"
	// BuildFileStateUploaded means the file is ready
	BuildFileStateUploaded BuildFileState = "uploaded"
	// BuildFileStateFailed means the file failed uploading
	BuildFileStateFailed BuildFileState = "failed"
)

// BuildFileType describes the type of a build file: patch, archive, signature, etc.
type BuildFileType string

const (
	// BuildFileTypePatch describes wharf patch files (.pwr)
	BuildFileTypePatch BuildFileType = "patch"
	// BuildFileTypeArchive describes canonical archive form (.zip)
	BuildFileTypeArchive BuildFileType = "archive"
	// BuildFileTypeSignature describes wharf signature files (.pws)
	BuildFileTypeSignature BuildFileType = "signature"
	// BuildFileTypeManifest is reserved
	BuildFileTypeManifest BuildFileType = "manifest"
	// BuildFileTypeUnpacked describes the single file that is in the build (if it was just a single file)
	BuildFileTypeUnpacked BuildFileType = "unpacked"
)

// BuildFileSubType describes the subtype of a build file: mostly its compression
// level. For example, rediff'd patches are "optimized", whereas initial patches are "default"
type BuildFileSubType string

const (
	// BuildFileSubTypeDefault describes default compression (rsync patches)
	BuildFileSubTypeDefault BuildFileSubType = "default"
	// BuildFileSubTypeGzip is reserved
	BuildFileSubTypeGzip BuildFileSubType = "gzip"
	// BuildFileSubTypeOptimized describes optimized compression (rediff'd / bsdiff patches)
	BuildFileSubTypeOptimized BuildFileSubType = "optimized"
)
