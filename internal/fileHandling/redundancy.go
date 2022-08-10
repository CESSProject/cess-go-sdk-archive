package fileHandling

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

func ReedSolomon_Restore(dir, fid string, datashards, rdushards int) error {
	outfn := filepath.Join(dir, fid)
	if rdushards == 0 {
		return os.Rename(outfn+".000", outfn)
	}
	if datashards+rdushards <= 6 {
		enc, err := reedsolomon.New(datashards, rdushards)
		if err != nil {
			return err
		}
		shards := make([][]byte, datashards+rdushards)
		for i := range shards {
			infn := fmt.Sprintf("%s.00%d", outfn, i)
			shards[i], err = ioutil.ReadFile(infn)
			if err != nil {
				shards[i] = nil
			}
		}

		// Verify the shards
		ok, _ := enc.Verify(shards)
		if !ok {
			err = enc.Reconstruct(shards)
			if err != nil {
				return err
			}
			ok, err = enc.Verify(shards)
			if !ok {
				return err
			}
		}
		f, err := os.Create(outfn)
		if err != nil {
			return err
		}

		err = enc.Join(f, shards, len(shards[0])*datashards)
		return err
	}

	enc, err := reedsolomon.NewStream(datashards, rdushards)
	if err != nil {
		return err
	}

	// Open the inputs
	shards, size, err := openInput(datashards, rdushards, outfn)
	if err != nil {
		return err
	}

	// Verify the shards
	ok, err := enc.Verify(shards)
	if !ok {
		shards, size, err = openInput(datashards, rdushards, outfn)
		if err != nil {
			return err
		}

		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				var outfn string
				if i < 10 {
					outfn = fmt.Sprintf("%s.00%d", outfn, i)
				} else {
					outfn = fmt.Sprintf("%s.0%d", outfn, i)
				}
				out[i], err = os.Create(outfn)
				if err != nil {
					return err
				}
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			return err
		}

		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				if err != nil {
					return err
				}
			}
		}
		shards, size, err = openInput(datashards, rdushards, outfn)
		ok, err = enc.Verify(shards)
		if !ok {
			return err
		}
		if err != nil {
			return err
		}
	}

	f, err := os.Create(outfn)
	if err != nil {
		return err
	}

	shards, size, err = openInput(datashards, rdushards, outfn)
	if err != nil {
		return err
	}

	err = enc.Join(f, shards, int64(datashards)*size)
	return err
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		var infn string
		if i < 10 {
			infn = fmt.Sprintf("%s.00%d", fname, i)
		} else {
			infn = fmt.Sprintf("%s.0%d", fname, i)
		}
		f, err := os.Open(infn)
		if err != nil {
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, 0, err
		}
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}
