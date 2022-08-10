package result

type UserSpaceDetails struct {
	TotalSpace     string
	UsedSpace      string
	RemainingSpace string
	Package_type   uint8
	Start          uint32
	Deadline       uint32
	State          string
}

type FileInfo struct {
	FileName  string `json:"file_name"`  //File name
	FileSize  int64  `json:"file_size"`  //File size
	FileState string `json:"file_state"` //File hash
}
type FindFileList struct {
	FileId string `json:"file_id"`
}
