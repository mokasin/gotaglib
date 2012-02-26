/*  Copyright 2012, mokasin
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package gotaglib

// #cgo LDFLAGS: -ltag_c
// #include <stdlib.h>
// #include <taglib/tag_c.h>
import "C"
import "unsafe"
import "fmt"

type TaggedFile struct {
	Filename string
	Title    string
	Artist   string
	Album    string
	Comment  string
	Genre    string
	Year     uint
	Track    uint
	Bitrate  uint
	Length   uint
}

type ErrUnreadableTag struct {}

func (e ErrUnreadableTag) Error() string {
	return fmt.Sprintf("Can not read tags of file.")
}

// Constructor of TaggedFile
func NewTaggedFile(filename string) (*TaggedFile, error) {

	//convert filename to a C-string and make sure to free it afterwards
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	taglib_file := C.taglib_file_new(cfilename)
	defer C.taglib_file_free((*C.TagLib_File)(taglib_file))
	if taglib_file == nil {
		return nil, fmt.Errorf("Unrecognized filetype or unreadable file.")
	}

	if C.taglib_file_is_valid(taglib_file) == 0 {
		return nil, ErrUnreadableTag{}
	}

	taglib_tag := C.taglib_file_tag(taglib_file)
	taglib_audioprop := C.taglib_file_audioproperties(taglib_file)

	return &TaggedFile{Filename: filename,
		Title:   C.GoString(C.taglib_tag_title(taglib_tag)),
		Artist:  C.GoString(C.taglib_tag_artist(taglib_tag)),
		Album:   C.GoString(C.taglib_tag_album(taglib_tag)),
		Comment: C.GoString(C.taglib_tag_comment(taglib_tag)),
		Genre:   C.GoString(C.taglib_tag_genre(taglib_tag)),
		Year:    uint(C.taglib_tag_year(taglib_tag)),
		Track:   uint(C.taglib_tag_track(taglib_tag)),
		Bitrate: uint(C.taglib_audioproperties_bitrate(taglib_audioprop)),
		Length:  uint(C.taglib_audioproperties_length(taglib_audioprop)),
	}, nil
}
