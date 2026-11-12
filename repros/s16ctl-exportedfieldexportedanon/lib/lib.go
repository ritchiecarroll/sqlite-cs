package lib

// S16 green control: the same shapes, with the inner alias EXPORTED.
type T__WIDL_name_00000003 = struct {
	FhInproc int32
}

type TuserHPALETTE = struct {
	FfContext int32
	Fu        T__WIDL_name_00000003
}

func Make(c, h int32) TuserHPALETTE {
	return TuserHPALETTE{FfContext: c, Fu: T__WIDL_name_00000003{FhInproc: h}}
}
