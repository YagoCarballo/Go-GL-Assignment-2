package wrapper

import (
	"time"
	"math"
	rng "github.com/leesper/go_rng"
)

type RandomMid struct {
	Height       [][]float32
	Magnitude    int32
	LandSize     int32
	Fractal      float32
	HeightFactor float32
	Scale        float32
	grng         *rng.GaussianGenerator
}

func NewRandomMid (height [][]float32, magnitude int32, fractal, heightFactor, scale float32) *RandomMid {
	return &RandomMid{
		height,
		magnitude,
		int32(math.Pow(2.0, float64(magnitude))) + 1,
		fractal,
		heightFactor,
		scale,
		rng.NewGaussianGenerator(time.Now().UnixNano()),
	}
}

func (data *RandomMid) CreateLandscape () {
	landMax := data.LandSize - 1
	scaleFactor	:= float32(1/(math.Pow(2.0, float64(data.Fractal / 2.0))))
	c := float32(math.Sqrt(1.0 - math.Pow(2.0, float64(2.0 * data.Fractal - 2.0))))
	scale := data.HeightFactor * float32(landMax) * scaleFactor * c * data.Scale;

	var avDist int32
	var x,y int32
	var nLevel, nFullLevel int32
	var xStep, yStep int32
	var iteration, yIteration, xIteration, yIter, xIter int32

	// Set the four corners of the array to different heights
	data.Height[0][0]				= 0 // dScale * (randomGauss());
	data.Height[landMax][0]			= 0 // dScale * (randomGauss());
	data.Height[0][landMax]			= 0 // dScale * (randomGauss());
	data.Height[landMax][landMax]	= 0 // dScale * (randomGauss());

	// Main Loop for each major iteration
	for iteration=1; iteration<=data.Magnitude; iteration++ {
		// Update scale factor
		//		if iteration != 1 {
		//			scale = scaleFactor * scale
		//		}

		scale = scaleFactor * scale;
		nLevel		= int32(math.Pow(2,float64(iteration-1)));
		nFullLevel	= int32(math.Pow(2,float64(iteration)));
		yStep		= landMax / nLevel;
		xStep		= landMax / nLevel;

		avDist = landMax / nFullLevel;
		y = landMax/nFullLevel;

		// Each centre point y iteration
		for yIteration=1; yIteration<=nLevel; yIteration ++ {
			x = landMax / nFullLevel;

			// For each X iteration
			for xIteration=1; xIteration<=nLevel; xIteration ++ {
				data.Height[x][y] = (data.Height[x-avDist][y-avDist] +
				data.Height[x-avDist][y+avDist] +
				data.Height[x+avDist][y-avDist] +
				data.Height[x+avDist][y+avDist]) / 4;
				data.Height[x][y] += (scale * float32(data.grng.StdGaussian()));

				x += xStep;
			}
			y += yStep;
		}

		// Update scale factor for next set of points
		scale = scale * scaleFactor;

		y = 0;
		yStep = landMax / nFullLevel;

		xStep = landMax / nLevel;

		// for each corner point y iteration
		for yIter=1; yIter<=(nFullLevel+1); yIter++ {
			// Split for even/odd
			if (yIter % 2) == 1 {
				// Odd Points
				x = landMax / nFullLevel;

				// For each even x iteration
				for xIter=1; xIter<=nLevel; xIter ++ {
					if y == 0 {
						data.Height[x][y] = (data.Height[x][y + avDist] +
						data.Height[x - avDist][y] +
						data.Height[x + avDist][y]) / 3;
					} else if y == landMax {
						data.Height[x][y] = (data.Height[x][y - avDist] +
						data.Height[x - avDist][y] +
						data.Height[x + avDist][y]) / 3;
					} else {
						data.Height[x][y] = (data.Height[x][y - avDist] +
						data.Height[x][y + avDist] +
						data.Height[x - avDist][y] +
						data.Height[x + avDist][y]) / 4;
					}

					data.Height[x][y] = data.Height[x][y] + scale * float32(data.grng.StdGaussian());
					x += xStep;
				}

				// {End of Even part of If}
			} else {
				// Split for Even
				x = 0;
				// For each even x iteration
				for xIter=1; xIter <= (int32(math.Pow(2, float64(iteration - 1))) + 1); xIter++ {
					if x == 0 {
						data.Height[x][y] = (data.Height[x][y - avDist] +
						data.Height[x][y + avDist] +
						data.Height[x + avDist][y]) / 3;
					} else if x == landMax {
						data.Height[x][y] = (data.Height[x][y - avDist] +
						data.Height[x][y + avDist] +
						data.Height[x - avDist][y]) / 3;
					} else {
						data.Height[x][y] = (data.Height[x][y - avDist] +
						data.Height[x][y + avDist] +
						data.Height[x - avDist][y] +
						data.Height[x + avDist][y]) / 4;
					}

					data.Height[x][y] = data.Height[x][y] + scale * float32(data.grng.StdGaussian());

					x += xStep;
				}
			}        // {End of odd part of If}
			y += yStep
		}

	} // End of main iteration loop
}

func (data *RandomMid) CreateLandscapeAdditional () {
	landMax := data.LandSize - 1
	scaleFactor	:= float32(1/(math.Pow(2.0, float64(data.Fractal / 2.0))))
	c := float32(math.Sqrt(1.0 - math.Pow(2.0, float64(2.0 * data.Fractal - 2.0))))
	scale := data.HeightFactor * float32(landMax) * scaleFactor * c * data.Scale

	var x,y int32
	var N, D, d int32
	var iteration int32
	nMaxLevel := int32(math.Pow(2, float64(data.Magnitude)))

	// Set the four corners of the array to different heights
	data.Height[0][0]				= scale * float32(data.grng.StdGaussian());
	data.Height[landMax][0]			= scale * float32(data.grng.StdGaussian());
	data.Height[0][landMax]			= scale * float32(data.grng.StdGaussian());
	data.Height[landMax][landMax]	= scale * float32(data.grng.StdGaussian());

	N = nMaxLevel;
	D = N;
	d = N / 2;

	// Main Loop for each major iteration
	for iteration=1; iteration<=data.Magnitude; iteration++ {
		// Update scale factor
		scale	*= scaleFactor;
		scaleFactor *= 1.01;

		for x=d; x<=(N-d); x+=D {
			for y=d; y<=(N-d); y+=D {
				data.Height[x][y] = ((data.Height[x+d][y+d] +
				data.Height[x+d][y-d] +
				data.Height[x-d][y+d] +
				data.Height[x-d][y-d]) / 4) + (scale * float32(data.grng.StdGaussian()))
			}
		}

		// Additional Random Displacements
		for x=0; x<=N; x+=D {
			for y = 0; y < N; y += D {
				data.Height[x][y] += (scale * float32(data.grng.StdGaussian()));
			}
		}

		// Update scale factor for next set of points
		scale = scale * scaleFactor;


		for x=d; x<=(N-d); x+=D {
			data.Height[x][0] = ((data.Height[x+d][0] +
			data.Height[x-d][0] +
			data.Height[x][d]) / 3) + (scale * float32(data.grng.StdGaussian()));
			data.Height[x][N] = ((data.Height[x+d][N] +
			data.Height[x-d][N] +
			data.Height[x][N-d]) / 3) + (scale * float32(data.grng.StdGaussian()));
			data.Height[0][x] = ((data.Height[0][x+d] +
			data.Height[0][x-d] +
			data.Height[d][x]) / 3) + (scale * float32(data.grng.StdGaussian()));
			data.Height[N][x] = ((data.Height[N][x+d] +
			data.Height[N][x-d] +
			data.Height[N-d][x]) / 3) + (scale * float32(data.grng.StdGaussian()));

		}

		// Interpolate and offset interior grid points
		for x:=d; x<=(N-d); x+=D {
			for y=D; y<=(N-d); y+=D {
				data.Height[x][y] = ((data.Height[x][y+d] +

				data.Height[x][y-d] +
				data.Height[x+d][y] +
				data.Height[x-d][y]) / 4) + (scale * float32(data.grng.StdGaussian()));
			}
		}
		for x=D; x<=(N-d); x+=D {
			for y=d; y<=(N-d); y+=D {
				data.Height[x][y] = ((data.Height[x][y+d] +
				data.Height[x][y-d] +
				data.Height[x+d][y] +
				data.Height[x-d][y]) / 4) + (scale * float32(data.grng.StdGaussian()));
			}
		}

		// Additional Random Displacements if required
		for x=0; x<=N; x += D {
			for y = 0; y <= N; y += D {
				data.Height[x][y] += (scale * float32(data.grng.StdGaussian()))
			}
		}

		for x=d; x<=(N-d); x += D {
			for y = d; y <= (N - d); y += D {
				data.Height[x][y] += (scale * float32(data.grng.StdGaussian()));
			}
		}

		D /= 2;
		d /= 2;

	}	// End of main iteration loop
}

func (data *RandomMid) SetMinMax () {
	var x, y int32
	var max, min float32 = 0.0, 0.0

	for y=0; y<data.LandSize; y++ {
		for x=0; x<data.LandSize; y++ {
			if data.Height[x][y] > max {
				max = data.Height[x][y]
			}

			if data.Height[x][y] < min {
				min = data.Height[x][y]
			}
		}
	}
}