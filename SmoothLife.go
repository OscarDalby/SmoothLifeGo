package main

import "gonum.org/v1/gonum/mat"

func ConstructSmoothLife(cm CellMath, mp Multipliers, br BasicRules, width int, height int) *SmoothLife {
	return &SmoothLife{
		width:  width,
		height: height,
		cm:     cm,
		mp:     mp,
		br:     br,
	}
}

type SmoothLife struct {
	width  int
	height int
	cm     CellMath
	mp     Multipliers
	br     BasicRules
	field  *mat.CDense
}

func (sl SmoothLife) Clear() {
	sl.field = mat.NewCDense(sl.height, sl.width, nil)
}

func (sl SmoothLife) Step() {
	var newField *mat.CDense = sl.cm.Fft2(sl.field)
}

//     def step(self):
//         """Do timestep and return field"

//         # To sum up neighbors, do kernel convolutions
//         # by multiplying in the frequency domain
//         # and converting back to spacial domain
//         field_ = np.fft.fft2(self.field)
//         M_buffer_ = field_ * self.multipliers.M
//         N_buffer_ = field_ * self.multipliers.N
//         M_buffer = np.real(np.fft.ifft2(M_buffer_))
//         N_buffer = np.real(np.fft.ifft2(N_buffer_))

//         # Apply transition rules
//         self.field = self.rules.s(N_buffer, M_buffer, self.field)
//         return self.field

//     def add_speckles(self, count=None, intensity=1):
//         """Populate field with random living squares

//         If count unspecified, do a moderately dense fill
//         """
//         if count is None:
//             count = int(
//                 self.width * self.height / ((self.multipliers.OUTER_RADIUS * 2) ** 2)
//             )
//         for i in range(count):
//             radius = int(self.multipliers.OUTER_RADIUS)
//             r = np.random.randint(0, self.height - radius)
//             c = np.random.randint(0, self.width - radius)
//             self.field[r : r + radius, c : c + radius] = intensity

// def show_animation():
//     w = 1 << 9
//     h = 1 << 9
//     # w = 1920
//     # h = 1080
//     sl = SmoothLife(h, w)
//     sl.add_speckles()
//     sl.step()

//     fig = plt.figure()
//     # Nice color maps: viridis, plasma, gray, binary, seismic, gnuplot
//     im = plt.imshow(
//         sl.field, animated=True, cmap=plt.get_cmap("viridis"), aspect="equal"
//     )

//     def animate(*args):
//         im.set_array(sl.step())
//         return (im,)

//     ani = animation.FuncAnimation(fig, animate, interval=60, blit=True)
//     plt.show()
