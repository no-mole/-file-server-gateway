package enum

import (
	"github.com/no-mole/neptune/enum"
)

var (
	ErrorGrpcClient enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4000,
		Msg:  "client get err",
	}

	ErrorSingleUpload enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4001,
		Msg:  "single upload err",
	}

	ErrorChunkUpload enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4002,
		Msg:  "chunk upload err",
	}

	ErrorNotFindConn enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4003,
		Msg:  "not find conn",
	}

	ErrorFileRead enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4004,
		Msg:  "read file err",
	}

	ErrorFileOpen enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4005,
		Msg:  "open file err",
	}

	ErrorDirOpen enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4006,
		Msg:  "open dir err",
	}

	ErrorRemoveFile enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4007,
		Msg:  "remove file err",
	}

	ErrorUploadFileBase64 enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4008,
		Msg:  "file base64 is nil",
	}

	ErrorGetFileMetadata enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4009,
		Msg:  "get file metadata err",
	}

	ErrorDownloadFile enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4010,
		Msg:  "download file err",
	}

	ErrorCreateFile enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4011,
		Msg:  "create file err",
	}

	ErrorWriteFile enum.ErrorNum = enum.ErrorNumEntry{
		Code: 4012,
		Msg:  "write file err",
	}
)
