package tartrans

import (
	"archive/tar"

	"go.polydawn.net/rio/fs"
	. "go.polydawn.net/rio/lib/errcat"
	"go.polydawn.net/timeless-api/rio"
)

// Mutate tar.Header fields to match the given fmeta.
func MetadataToTarHdr(fmeta *fs.Metadata, hdr *tar.Header) {
	// TODO
}

// Mutate fs.Metadata fields to match the given tar header.
// Does not check for names that go above '.'; caller may want to do that.
func TarHdrToMetadata(hdr *tar.Header, fmeta *fs.Metadata) error {
	fmeta.Name = fs.MustRelPath(hdr.Name) // FIXME should not use the 'must' path
	fmeta.Type = tarTypeToFsType(hdr.Typeflag)
	if fmeta.Type == fs.Type_Invalid {
		return Errorf(rio.ErrWareCorrupt, "corrupt tar: %q is not a known file type", hdr.Typeflag)
	}
	fmeta.Perms = fs.Perms(hdr.Mode & 07777)
	fmeta.Uid = uint32(hdr.Uid)
	fmeta.Gid = uint32(hdr.Gid)
	fmeta.Size = hdr.Size
	fmeta.Linkname = hdr.Linkname
	fmeta.Devmajor = hdr.Devmajor
	fmeta.Devminor = hdr.Devminor
	fmeta.Mtime = hdr.ModTime
	fmeta.Xattrs = hdr.Xattrs
	return nil
}

func tarTypeToFsType(tarType byte) fs.Type {
	switch tarType {
	case tar.TypeReg, tar.TypeRegA:
		return fs.Type_File
	case tar.TypeLink:
		return fs.Type_Hardlink
	case tar.TypeSymlink:
		return fs.Type_Symlink
	case tar.TypeChar:
		return fs.Type_CharDevice
	case tar.TypeBlock:
		return fs.Type_Device
	case tar.TypeDir:
		return fs.Type_Dir
	case tar.TypeFifo:
		return fs.Type_NamedPipe
	// Notice that tar does not have a type for socket files
	default:
		return fs.Type_Invalid
	}
}
