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
	FileName    string `json:"file_name"`   //File name
	FileSize    int64  `json:"file_size"`   //File size
	FileHash    string `json:"file_hash"`   //File hash
	Public      bool   `json:"public"`      //Public or not
	Backups     int8   `json:"backups"`     //Number of backups
	Downloadfee int64  `json:"downloadfee"` //Download fee
}
type FindFileList struct {
	FileId string `json:"file_id"`
}
