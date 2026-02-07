package krun

// LogLevel represents the verbosity level for logging.
type LogLevel uint32

const (
	LogLevelOff   LogLevel = 0
	LogLevelError LogLevel = 1
	LogLevelWarn  LogLevel = 2
	LogLevelInfo  LogLevel = 3
	LogLevelDebug LogLevel = 4
	LogLevelTrace LogLevel = 5
)

// LogStyle controls terminal escape sequence usage in log output.
type LogStyle uint32

const (
	LogStyleAuto   LogStyle = 0
	LogStyleAlways LogStyle = 1
	LogStyleNever  LogStyle = 2
)

// LogTargetDefault uses the default log target (stderr).
const LogTargetDefault = -1

// LogOptionNoEnv disallows environment variables from overriding log settings.
const LogOptionNoEnv uint32 = 1

// DiskFormat represents the format of a disk image.
type DiskFormat uint32

const (
	DiskFormatRaw   DiskFormat = 0
	DiskFormatQcow2 DiskFormat = 1
	// DiskFormatVmdk only supports FLAT/ZERO formats without delta links.
	DiskFormatVmdk DiskFormat = 2
)

// SyncMode controls VIRTIO_BLK_F_FLUSH behavior.
type SyncMode uint32

const (
	// SyncNone ignores VIRTIO_BLK_F_FLUSH. WARNING: may lead to data loss.
	SyncNone SyncMode = 0
	// SyncRelaxed honors flush requests but relaxes strict hardware syncing on macOS.
	// This is the recommended mode.
	SyncRelaxed SyncMode = 1
	// SyncFull strictly flushes buffers to physical disk.
	SyncFull SyncMode = 2
)

// KernelFormat represents the format of a kernel image.
type KernelFormat uint32

const (
	KernelFormatRaw       KernelFormat = 0
	KernelFormatELF       KernelFormat = 1
	KernelFormatPEGz      KernelFormat = 2
	KernelFormatImageBz2  KernelFormat = 3
	KernelFormatImageGz   KernelFormat = 4
	KernelFormatImageZstd KernelFormat = 5
)

// Feature represents a compile-time feature flag for [HasFeature].
type Feature uint64

const (
	FeatureNet              Feature = 0
	FeatureBLK              Feature = 1
	FeatureGPU              Feature = 2
	FeatureSND              Feature = 3
	FeatureInput            Feature = 4
	FeatureEFI              Feature = 5
	FeatureTEE              Feature = 6
	FeatureAMDSEV           Feature = 7
	FeatureIntelTDX         Feature = 8
	FeatureAWSNitro         Feature = 9
	FeatureVirglResourceMap2 Feature = 10
)

// Network flags.
const (
	// NetFlagVfkit sends the VFKIT magic after establishing the connection,
	// as required by gvproxy in vfkit mode.
	NetFlagVfkit uint32 = 1 << 0
)

// Network feature flags (from virtio_net.h).
const (
	NetFeatureCsum      uint32 = 1 << 0
	NetFeatureGuestCsum uint32 = 1 << 1
	NetFeatureGuestTSO4 uint32 = 1 << 7
	NetFeatureGuestTSO6 uint32 = 1 << 8
	NetFeatureGuestUFO  uint32 = 1 << 10
	NetFeatureHostTSO4  uint32 = 1 << 11
	NetFeatureHostTSO6  uint32 = 1 << 12
	NetFeatureHostUFO   uint32 = 1 << 14
)

// CompatNetFeatures is the set of features enabled by the deprecated
// krun_set_passt_fd and krun_set_gvproxy_path functions.
const CompatNetFeatures = NetFeatureCsum | NetFeatureGuestCsum |
	NetFeatureGuestTSO4 | NetFeatureGuestUFO |
	NetFeatureHostTSO4 | NetFeatureHostUFO

// TSI (Transparent Socket Impersonation) feature flags for vsock.
const (
	TSIHijackInet uint32 = 1 << 0
	TSIHijackUnix uint32 = 1 << 1
)

// Virglrenderer flags for GPU configuration.
const (
	VirglUseEGL           uint32 = 1 << 0
	VirglThreadSync       uint32 = 1 << 1
	VirglUseGLX           uint32 = 1 << 2
	VirglUseSurfaceless   uint32 = 1 << 3
	VirglUseGLES          uint32 = 1 << 4
	VirglUseExternalBlob  uint32 = 1 << 5
	VirglVenus            uint32 = 1 << 6
	VirglNoVirgl          uint32 = 1 << 7
	VirglUseAsyncFenceCB  uint32 = 1 << 8
	VirglRenderServer     uint32 = 1 << 9
	VirglDRM              uint32 = 1 << 10
)

// MaxDisplays is the maximum number of displays (same as VIRTIO_GPU_MAX_SCANOUTS).
const MaxDisplays = 16
