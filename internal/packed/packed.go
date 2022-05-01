package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBg4GBobZ0XwoAEhBg4GZLz89Iy0/UhlF5lYm5OaAgrA2N0em5SyFnHvlYFnmMf9y7P2CU/82iK2B02wY2dS/mmdL7/pn/r+SHRk/r/+90+MTuo2co3LOWqPhBamn9vp8+ub5k9h/qcXsvIxM3sesoQ32+xNN36iWre9daVl81ur2Dxu3dtx1b734ZXdtyTqPj31zBv/bLOvVtDfiWeMapPOqCkOvOYzva6+A/zs5Q11RepVgdYNy+evUJimXP2BtVH1gwM//8HeLNzCHddjpjJwMAgyMjAAPMhA8NyNB+ywX0I9tXs9NwkkGZkJQHejEwizIgAQjYYFEAwsKQRROIJLoRB2N0BAQIM/x0fwQ1CchUrG0iaiYGJoY2BgUGeEcQDBAAA//+s4s9LuwEAAA=="); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
