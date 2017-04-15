package tartrans

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"strings"

	"go.polydawn.net/rio"
	"go.polydawn.net/rio/fs"
)

func Extract(
	ctx context.Context,
	destBasePath fs.AbsolutePath,
	filters rio.Filters,
	tr *tar.Reader,
) rio.Error {
	// Iterate over each tar entry, mutating filesystem as we go.
	for {
		fmeta := fs.Metadata{}
		thdr, err := tr.Next()

		// Check for done.
		if err == io.EOF {
			return nil // sucess!  end of archive.
		}
		if err != nil {
			return rio.ErrWareCorrupt{
				Msg: fmt.Sprintf("corrupt tar: %s", err),
			}
		}
		if ctx.Err() != nil {
			return rio.Cancelled{}
		}

		// Reshuffle metainfo to our default format.
		if err := TarHdrToMetadata(thdr, &fmeta); err != nil {
			return rio.ErrWareCorrupt{
				Msg: fmt.Sprintf("corrupt tar: %s", err),
			}
		}
		if strings.HasPrefix(fmeta.Name.String(), "..") {
			return rio.ErrWareCorrupt{
				Msg: "corrupt tar: paths that use '../' to leave the base dir are invalid",
			}
		}

		// Apply filters.
		ApplyMaterializeFilters(&fmeta, filters)
	}
	return nil
}

// Mutate tar.Header fields to match the given fmeta.
func MetadataToTarHdr(fmeta *fs.Metadata, hdr *tar.Header) {

}

// Mutate fs.Metadata fields to match the given tar header.
// Does not check for names that go above '.'; caller may want to do that.
func TarHdrToMetadata(hdr *tar.Header, fmeta *fs.Metadata) error {
	fmeta.Name = fs.MustRelPath(hdr.Name) // FIXME should not use the 'must' path
	fmeta.Mode = hdr.FileInfo().Mode()
	fmeta.Uid = hdr.Uid
	fmeta.Gid = hdr.Gid
	fmeta.Size = hdr.Size
	fmeta.Linkname = hdr.Linkname
	fmeta.Devmajor = hdr.Devmajor
	fmeta.Devminor = hdr.Devminor
	fmeta.Mtime = hdr.ModTime
	fmeta.Xattrs = hdr.Xattrs
	return nil
}
