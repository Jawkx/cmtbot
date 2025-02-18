package model

type FileStatus int

const (
	StatusAdded FileStatus = iota
	StatusModified
	StatusDeleted
	StatusRenamed
	StatusCopied
	StatusUnmodified
	StatusUnknown
)

type StagedFile struct {
	Path   string
	Status FileStatus
}

func (s FileStatus) Name() string {
	switch s {
	case StatusAdded:
		return "Added"
	case StatusModified:
		return "Modified"
	case StatusDeleted:
		return "Deleted"
	case StatusRenamed:
		return "Renamed"
	case StatusCopied:
		return "Copied"
	case StatusUnmodified:
		return "Unmodified"
	case StatusUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

func (s FileStatus) String() string {
	switch s {
	case StatusAdded:
		return "A"
	case StatusModified:
		return "M"
	case StatusDeleted:
		return "D"
	case StatusRenamed:
		return "R"
	case StatusCopied:
		return "C"
	case StatusUnmodified:
		return "U"
	default:
		return "?"
	}
}
