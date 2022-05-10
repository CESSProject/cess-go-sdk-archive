package result

type UserHoldSpaceDetails struct {
	PurchasedSpace string `json:"purchased_space"`
	UsedSpace      string `json:"used_space"`
	RemainingSpace string `json:"remaining_space"`
}

type FileInfo struct {
	File_Name   string `json:"file_name"`   //File name
	FileSize    int64  `json:"file_size"`   //File size
	FileHash    string `json:"file_hash"`   //File hash
	Public      bool   `json:"public"`      //Public or not
	Backups     int8   `json:"backups"`     //Number of backups
	Downloadfee int64  `json:"downloadfee"` //Download fee
}
