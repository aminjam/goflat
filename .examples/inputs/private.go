package main

type Private struct {
	Password string
	Secret   string
	Key      string
}

func NewPrivate() Private {
	return Private{
		Password: "team3",
		Secret:   "cloud-foundry",
		Key: `-----BEGIN RSA PRIVATE KEY-----
P1231231313131231231312T2iBZ1jrI6ggQAHH3/LQNd1yN5oa8cE/kS9ll923u
4ksX0oo0enasdhakjhdakjdshakjshdkajshdakjdhakjdhHz9jtJ8t6165OfIu+
ApTtie8nAXK9asdhasjdhaksdwequwhequheqiuwheqiuhjsIelc11U4bQN+NKuv
Dg3asdadawuhuiqeuqihjsdhaksjhdajdsasjdhasjdhaasdz1VI9E3puhubTIVg
b12Px14IO23By2bd+jn+kgKLDS1Hwo8mfcOPbKx8gOxDzXe1o9/59W9dBtSq9ytT
DR6LCaPye6sjQIAXDHOQOCqzREBKceO1dy54LwIDAQABAoIBA1HWQk1a1yovex1X
vHna0tQxALQXeqxXCl3J1vdJDi12Eq51e5fEyP/BRN2Qfcsp944d4C5DZJU54Ite
/rNH1AWrfcJIi4HReWDIagqY1l9uZVl1ABQSjahdjakhsdjashdaaqweqweweqql
wfZna3s0dxcYLXACPY5Pn8ekeQXyJBxT1Ys51ASDGASDQasdashdjkwADSS111iK
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
UK4+r6EBasdasdasdasdsadsadasdqweqweuqweihukuKkT90nrwQAWXId1t9LPy
1b8cY/YuwKINO98PE/gvzOcsp1NN3ucWOm+0Bo0m151fzHDVKQB6oPVNEEa1a05Q
wfZna3s0dxcYLXACPY5Pn8ekeQXyJBxT1Ys51ASDGASDQasdashdjkwADSS111iK
NRmT+e/awtpqH+py0L5Pb4jP9A0o634dt/sAr+1IJkip689qEHPQlUBHqBHSPQoH
WP9n9diasdadasdasdoiqweioqwuoiQWEQWj4zORs6i8maxzzaq115iTAdUxn0PC
I+5NsEAPsqzwfCAv3l5DiZjiZHSqjiyZUZRkSfEBgYBhRrPJtb4qrKlkwLKhYP59
15KCBv+sBmd5X8cczKHi9fpdg/tTVgVfNOlCiIfkDUUk8P4v1WyfB5V9xDJbnq/h
b2kClRw5PeT/X59WjWR9tBT4lmPWYrc0UTT5HWOEue4NjxHDHBfsaf6nntrEDTgv
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
y848Dy3R89IWfWhcvtJlBiPu3213123123ASDASDASDASDasdascvbahjweqqK2d
eqz12k9Wu+d31C9t8UE5Ng1vPfLQZs6mx9R08vwzqlNezaPLSS+m
-----END RSA PRIVATE KEY-----`,
	}
}
