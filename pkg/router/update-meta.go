package router

type UpdateMeta struct {
	PassThrough bool
}

func PassThrough() *UpdateMeta {
	return &UpdateMeta{PassThrough: true}
}
