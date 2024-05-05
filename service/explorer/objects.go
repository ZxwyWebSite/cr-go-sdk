package explorer

// ItemMoveService 处理多文件/目录移动
type ItemMoveService struct {
	SrcDir string        `json:"src_dir" binding:"required,min=1,max=65535"`
	Src    ItemIDService `json:"src"`
	Dst    string        `json:"dst" binding:"required,min=1,max=65535"`
}

// ItemRenameService 处理多文件/目录重命名
type ItemRenameService struct {
	Src     ItemIDService `json:"src"`
	NewName string        `json:"new_name" binding:"required,min=1,max=255"`
}

// ItemService 处理多文件/目录相关服务
type ItemService struct {
	Items []uint `json:"items"`
	Dirs  []uint `json:"dirs"`
}

// ItemIDService 处理多文件/目录相关服务，字段值为HashID，可通过Raw()方法获取原始ID
type ItemIDService struct {
	Items []string `json:"items"`
	Dirs  []string `json:"dirs"`
	// Source     *ItemService
	Force      *bool `json:"force,omitempty"`
	UnlinkOnly *bool `json:"unlink,omitempty"`
}

// ItemCompressService 文件压缩任务服务
type ItemCompressService struct {
	Src  ItemIDService `json:"src"`
	Dst  string        `json:"dst" binding:"required,min=1,max=65535"`
	Name string        `json:"name" binding:"required,min=1,max=255"`
}

// ItemRelocateService 文件转移任务服务
type ItemRelocateService struct {
	Src         ItemIDService `json:"src"`
	DstPolicyID string        `json:"dst_policy_id" binding:"required"`
}

// ItemDecompressService 文件解压缩任务服务
type ItemDecompressService struct {
	Src      string `json:"src"`
	Dst      string `json:"dst" binding:"required,min=1,max=65535"`
	Encoding string `json:"encoding"`
}

// ItemPropertyService 获取对象属性服务
type ItemPropertyService struct {
	ID        string `binding:"required"`
	TraceRoot bool   `form:"trace_root"`
	IsFolder  bool   `form:"is_folder"`
}
